package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/bizerr"
	"github.com/ragpanda/go-toolkit/log"
	"github.com/ragpanda/go-toolkit/utils"
	"github.com/ragpanda/go-toolkit/utils/ratelimit"
)

type HttpClient struct {
	client        http.Client
	defaultOption HttpOptionalArgs
	rateLimit     map[string]*ratelimit.RateLimit
	lock          sync.RWMutex
}

func NewHttpClient(args HttpOptionalArgs) *HttpClient {
	if args.MaxConnectionPerHost == 0 {
		args.MaxConnectionPerHost = 100
	}
	if args.GlobalTimeoutSec == 0 {
		args.GlobalTimeoutSec = 15
	}
	if args.MaxIdleConnectionPerHost == 0 {
		args.MaxIdleConnectionPerHost = 30
	}
	if args.IdleConnectionTimoutSec == 0 {
		args.IdleConnectionTimoutSec = 30
	}

	c := &HttpClient{
		client: http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:    args.MaxIdleConnectionPerHost,
				MaxConnsPerHost:        args.MaxConnectionPerHost,
				IdleConnTimeout:        time.Duration(args.IdleConnectionTimoutSec) * time.Second,
				ResponseHeaderTimeout:  0,
				ExpectContinueTimeout:  0,
				TLSNextProto:           nil,
				ProxyConnectHeader:     nil,
				GetProxyConnectHeader:  nil,
				MaxResponseHeaderBytes: 0,
				WriteBufferSize:        0,
				ReadBufferSize:         0,
				ForceAttemptHTTP2:      false,
			},
			Timeout: time.Duration(args.GlobalTimeoutSec) * time.Second,
		},
		defaultOption: args,
	}

	return c
}

type HttpOptionalArgs struct {
	QPSLimit *int64

	Header                 map[string]string
	Method                 string
	UrlValues              url.Values
	TimeoutSec             int64
	DisableCheckStatusCode bool

	MaxRetryTimes        int
	RetryIntervalMillSec int64

	GlobalHttpConfig
}

type GlobalHttpConfig struct {
	DefaultHost              string
	MaxConnectionPerHost     int
	MaxIdleConnectionPerHost int
	IdleConnectionTimoutSec  int64
	ForceTryHttp2            bool

	GlobalTimeoutSec int64
}
type HttpOption func(*HttpOptionalArgs)

func (self *HttpClient) DoJson(ctx context.Context, urlStr string, body interface{}, httpOptions ...HttpOption) *HttpResultSet {
	option := self.getDefaultOption()
	for _, opt := range httpOptions {
		opt(option)
	}
	result := &HttpResultSet{}

	if !strings.HasPrefix(urlStr, "http") {
		if option.DefaultHost == "" {
			result.err = bizerr.ErrInternalError.WithMessage("url host/default host is empty")
			return result
		}

		if !strings.HasSuffix(option.DefaultHost, "/") && !strings.HasPrefix(urlStr, "/") {
			urlStr = option.DefaultHost + "/" + urlStr
		} else if strings.HasSuffix(option.DefaultHost, "/") && strings.HasPrefix(urlStr, "/") {
			urlStr = option.DefaultHost + urlStr[1:]
		} else {
			urlStr = option.DefaultHost + urlStr
		}
	}

	urlObj, err := url.Parse(urlStr)
	if err != nil {
		log.Error(ctx, "parse url failed %s", err.Error())
		result.err = err
		return result
	}

	if option.UrlValues != nil {
		queryParams := urlObj.Query()
		for k, vs := range option.UrlValues {
			queryParams.Set(k, vs[0])
		}
		urlObj.RawQuery = queryParams.Encode()
	}

	if option.Method == "" {
		option.Method = http.MethodPost
	}

	if option.RetryIntervalMillSec == 0 {
		option.RetryIntervalMillSec = 100
	}

	var rateLimit *ratelimit.RateLimit
	if option.QPSLimit != nil {
		self.lock.RLock()
		rateLimit = self.rateLimit[urlObj.Host]
		if rateLimit == nil {
			self.lock.RUnlock()
			self.lock.Lock()
			rateLimit = ratelimit.NewRateLimit(fmt.Sprintf("http-%s", urlObj.Host), int(*option.QPSLimit), 1*time.Second)
			self.rateLimit[urlObj.Host] = rateLimit
			self.lock.Unlock()
		} else {
			self.lock.RUnlock()
		}
	}

	httpReq := &http.Request{
		Method: option.Method,
		URL:    urlObj,
		Header: make(http.Header),
		Body:   io.NopCloser(bytes.NewBuffer(utils.MustJsonEncodeBytes(body))),
	}

	for k, v := range option.Header {
		httpReq.Header.Set(k, v)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if option.TimeoutSec != 0 {
		newCtx, cancel := context.WithTimeout(ctx, time.Duration(option.TimeoutSec)*time.Second)
		defer cancel()
		ctx = newCtx
	}

	result.Url = urlStr
	err = utils.Retry(ctx, func(ctx context.Context) error {
		if rateLimit != nil {
			rateLimit.Take(ctx)
		}
		resp, err := self.client.Do(httpReq)
		if err != nil {
			log.Warn(ctx, "http request error %s", err.Error())
			return err
		}

		if resp.Body == nil {
			log.Warn(ctx, "http response body is nil")
			return bizerr.ErrInternalError.WithMessage("http response body is nil")
		}
		defer resp.Body.Close()

		result.StatusCode = resp.StatusCode
		if !option.DisableCheckStatusCode && (resp.StatusCode < 200 || resp.StatusCode > 226) {
			log.Warn(ctx, "http request status code error %s, %d", urlStr, resp.StatusCode)
			return bizerr.ErrInternalError.WithMessage(
				fmt.Sprintf("http response code invalid, url=`%s`, status_code=`%d`", urlStr, resp.StatusCode))
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Warn(ctx, "http read body error %s", err.Error())
			return err
		}

		result.Body = respBody
		return nil
	}, utils.RetryMaxTimes(option.MaxRetryTimes, time.Duration(option.RetryIntervalMillSec)*time.Millisecond))

	if err != nil {
		result.err = err
		log.Warn(ctx, "http request failed %s", err.Error())
		return result
	}

	return result
}

func (self *HttpClient) getDefaultOption() *HttpOptionalArgs {
	data := utils.JsonDeepCopy(self.defaultOption)
	return data
}

type HttpResultSet struct {
	Url        string
	StatusCode int
	Body       []byte

	err error
}

func (self *HttpResultSet) Error() error {
	return self.err
}

func (self *HttpResultSet) Unmarshal(ctx context.Context, v interface{}) error {
	if self.err != nil {
		return self.err
	}

	err := utils.Unmarshal(self.Body, v)
	if err != nil {
		log.Warn(ctx, "http unmarshal error %s, body=`%s`", err.Error(), self.Body)
		return err
	}

	return nil
}
