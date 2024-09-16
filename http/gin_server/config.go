package gin_server

type GinConfig struct {
	Addr string `yaml:"Addr" json:"Addr"`
	Mode string `yaml:"Mode" json:"Mode"`

	EnableBaseMw bool `yaml:"EnableBaseMw" json:"EnableBaseMw"`
	EnablePprof  bool `yaml:"EnablePprof" json:"EnablePprof"`

	ProfilePath string           `yaml:"ProfilePath" json:"ProfilePath"`
	CORS        *CORSConfig      `yaml:"CORS" json:"CORS"`
	RateLimit   *RateLimitConfig `yaml:"RateLimit" json:"RateLimit"`

	GracefulExitSec int64 `yaml:"GracefulExitSec" json:"GracefulExitSec"`
}

type CORSConfig struct {
	Enable       bool     `yaml:"Enable" json:"Enable"`
	AllowOrigins []string `yaml:"AllowOrigins" json:"AllowOrigins"`
}
