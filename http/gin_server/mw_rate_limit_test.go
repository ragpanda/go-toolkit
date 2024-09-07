package gin_server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinRateLimiter(t *testing.T) {
	// 创建测试用的Gin引擎
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 创建限流配置
	globalLimit := 5
	userLimit := 2
	ipLimit := 3
	config := &RateLimitConfig{
		Rules: []*RateLimitRule{
			{
				Mode:            ModeTokenBucket,
				MatchPathPrefix: "/api/",
				// PerUserLimit:    &userLimit,
				// PerIPLimit:      &ipLimit,
				GlobalLimit:  &globalLimit,
				CycleSecond:  1,
				BreakIfMatch: true,
			},
		},
	}

	// 创建限流中间件
	limiter := NewGinRateLimiter(config)
	router.Use(limiter.RateLimitMW())

	// 注册测试路由
	router.GET("/api/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 测试用例1：全局限流
	t.Run("Global Limit", func(t *testing.T) {
		limiter.Reload(config)
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if i < globalLimit {
				assert.Equal(t, http.StatusOK, w.Code, i)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, w.Code, i)
			}
		}
	})

	// 测试用例2：用户限流
	config = &RateLimitConfig{
		Rules: []*RateLimitRule{
			{
				Mode:            ModeTokenBucket,
				MatchPathPrefix: "/api/",
				PerUserLimit:    &userLimit,
				CycleSecond:     1,
				BreakIfMatch:    true,
			},
		},
	}
	t.Run("User Limit", func(t *testing.T) {
		limiter.Reload(config)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			req.Header.Set("X-User-ID", "user1")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if i < userLimit {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, w.Code)
			}
		}
	})

	// 测试用例3：IP限流

	config = &RateLimitConfig{
		Rules: []*RateLimitRule{
			{
				Mode:            ModeTokenBucket,
				MatchPathPrefix: "/api/",
				PerIPLimit:      &ipLimit,
				CycleSecond:     1,
				BreakIfMatch:    true,
			},
		},
	}
	t.Run("IP Limit", func(t *testing.T) {
		limiter.Reload(config)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			req.RemoteAddr = "127.0.0.1:8080"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if i < ipLimit {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, w.Code)
			}
		}
	})

	// 测试用例4：合法请求
	config = &RateLimitConfig{
		Rules: []*RateLimitRule{
			{
				Mode:            ModeTokenBucket,
				MatchPathPrefix: "/api/",
				GlobalLimit:     &globalLimit,
				CycleSecond:     1,
				BreakIfMatch:    true,
			},
		},
	}

	t.Run("Legal Request", func(t *testing.T) {
		limiter.Reload(config)
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-User-ID", "user2")
		req.RemoteAddr = "127.0.0.2:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试用例5：限流规则切换
	config = &RateLimitConfig{
		Rules: []*RateLimitRule{
			{
				Mode:            ModeTokenBucket,
				MatchPathPrefix: "/api/",
				GlobalLimit:     &globalLimit,
				CycleSecond:     1,
				BreakIfMatch:    true,
			},
		},
	}
	t.Run("Rule Switch", func(t *testing.T) {
		newGlobalLimit := 1
		newUserLimit := 1
		newIPLimit := 1
		newConfig := &RateLimitConfig{
			Rules: []*RateLimitRule{
				{
					Mode:            ModeTokenBucket,
					MatchPathPrefix: "/api/",
					PerUserLimit:    &newUserLimit,
					PerIPLimit:      &newIPLimit,
					GlobalLimit:     &newGlobalLimit,
					CycleSecond:     1,
					BreakIfMatch:    true,
				},
			},
		}
		limiter.Reload(newConfig)

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-User-ID", "user3")
		req.RemoteAddr = "127.0.0.3:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		req = httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-User-ID", "user3")
		req.RemoteAddr = "127.0.0.3:8080"
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	// 测试用例6：无限流规则
	t.Run("No Limit", func(t *testing.T) {
		newConfig := &RateLimitConfig{
			Rules: []*RateLimitRule{
				{
					Mode:            ModeTokenBucket,
					MatchPathPrefix: "/api/",
					CycleSecond:     1,
					BreakIfMatch:    true,
				},
			},
		}
		limiter.Reload(newConfig)

		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			req.Header.Set("X-User-ID", "user4")
			req.RemoteAddr = "127.0.0.4:8080"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
