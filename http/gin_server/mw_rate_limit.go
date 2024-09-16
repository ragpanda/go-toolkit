package gin_server

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/biz"
	"github.com/ragpanda/go-toolkit/log"
	"golang.org/x/time/rate"
)

type GinRateLimiter struct {
	RuleItems []*RuleItem
}

type RateLimitConfig struct {
	Rules []*RateLimitRule `yaml:"Rules" json:"Rules"`
}

type RateLimitRule struct {
	// Mode 限流模式
	Mode RateLimitRuleMode `yaml:"Mode" json:"Mode"`
	// MatchPathPrefix 匹配路径前缀
	MatchPathPrefix string `yaml:"MatchPathPrefix" json:"MatchPathPrefix"`

	// PerUserLimit 用户限流值
	PerUserLimit *int `yaml:"PerUserLimit" json:"PerUserLimit"`
	// PerIPLimit IP 限流值
	PerIPLimit *int `yaml:"PerIPLimit" json:"PerIPLimit"`
	// GlobalLimit 全局限流值
	GlobalLimit *int `yaml:"GlobalLimit" json:"GlobalLimit"`

	// CycleSecond 限流周期，单位秒
	CycleSecond int `yaml:"CycleSecond" json:"CycleSecond"`

	// BreakIfMatch 为 true 时，匹配到该规则后不再继续匹配后续规则
	BreakIfMatch bool `yaml:"BreakIfMatch" json:"BreakIfMatch"`
}

type RuleItem struct {
	config        RateLimitRule
	globalLimiter *rate.Limiter
	userLimiters  sync.Map
	ipLimiters    sync.Map
}

type RateLimitRuleMode string

const (
	ModeTokenBucket RateLimitRuleMode = "token_bucket"
	ModeLeakyBucket RateLimitRuleMode = "leak_bucket"
)

func NewGinRateLimiter(config *RateLimitConfig) *GinRateLimiter {
	limiter := &GinRateLimiter{}
	limiter.Reload(config)
	return limiter
}

func (l *GinRateLimiter) RateLimitMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		b := biz.GetBizData(c)
		var userID, fromIP string
		if b == nil {
			userID = c.Request.Header.Get("X-User-ID")
			fromIP = c.RemoteIP()
		} else {
			userID = b.UserID
			fromIP = b.FromIP
		}
		if !l.Check(c, c.Request.URL.Path, userID, fromIP) {
			c.AbortWithStatus(429) // Too Many Requests
			return
		}

		c.Next()
	}
}
func (l *GinRateLimiter) Check(c context.Context, path, userID, fromIP string) bool {
	for _, item := range l.RuleItems {
		if item.match(path) {
			if !item.allowGlobal() || !item.allowUser(userID) || !item.allowIP(fromIP) {
				return false
			}

			if item.config.BreakIfMatch {
				break
			}

		}
	}
	return true
}

func (l *GinRateLimiter) Reload(config *RateLimitConfig) {
	var items []*RuleItem
	for _, rule := range config.Rules {
		item := &RuleItem{
			config:        *rule,
			globalLimiter: newLimiter(rule.Mode, rule.GlobalLimit, rule.CycleSecond),
		}
		items = append(items, item)
	}

	l.RuleItems = items
}

func (item *RuleItem) match(path string) bool {
	return len(path) >= len(item.config.MatchPathPrefix) && path[:len(item.config.MatchPathPrefix)] == item.config.MatchPathPrefix
}

func (item *RuleItem) allowGlobal() bool {
	if item.config.GlobalLimit == nil {
		return true
	}

	return item.globalLimiter.Allow()
}

func (item *RuleItem) allowUser(user string) bool {
	if item.config.PerUserLimit == nil {
		return true
	}

	limiter, _ := item.userLimiters.LoadOrStore(user, newLimiter(item.config.Mode, item.config.PerUserLimit, item.config.CycleSecond))
	if limiter == nil {
		return true
	}
	return limiter.(*rate.Limiter).Allow()
}

func (item *RuleItem) allowIP(ip string) bool {
	if item.config.PerIPLimit == nil {
		return true
	}

	limiter, _ := item.ipLimiters.LoadOrStore(ip, newLimiter(item.config.Mode, item.config.PerIPLimit, item.config.CycleSecond))
	if limiter == nil {
		return true
	}
	return limiter.(*rate.Limiter).Allow()
}

func newLimiter(mode RateLimitRuleMode, limit *int, cycleSec int) *rate.Limiter {
	if limit == nil || cycleSec <= 0 {
		return nil
	}

	switch mode {
	case ModeTokenBucket:
		return rate.NewLimiter(rate.Limit(float64(*limit)/float64(cycleSec)), *limit)
	case ModeLeakyBucket, "":
		return rate.NewLimiter(rate.Every(time.Duration(cycleSec)*time.Second), *limit)
	default:
		log.Error(context.Background(), "unsupported rate limit mode: %s", mode)
		return nil
	}
}
