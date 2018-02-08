package toolez

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/vettu/allez/core"
	"github.com/vettu/allez/toolez/keygen"
)

func New() core.ISiteInfo {
	s := &toolez{
		siteName: "toolez",
	}
	return s
}

type toolez struct {
	siteName   string
	rootConfig *core.Config
	config     *Config
	router     *gin.Engine

	middlewares []gin.HandlerFunc
}

func (tz *toolez) SiteFile(paths ...string) string {
	paths = append([]string{tz.rootConfig.RootPath}, paths...)
	return filepath.Join(paths...)
}

func (tz *toolez) Init(config *core.Config) {
	tz.rootConfig = config
	c := &Config{}
	mapstructure.Decode(config.Site[tz.siteName], c)
	config.Site[tz.siteName] = c
	tz.config = c

	tz.router = gin.New()
}

func (tz *toolez) DelayUse(middleware ...gin.HandlerFunc) core.ISiteInfo {
	tz.middlewares = append(tz.middlewares, middleware...)
	return tz
}
func (tz *toolez) Use() {
	tz.router.Use(tz.middlewares...)
}

func (tz *toolez) RegRouter() {
	r := tz.router
	r.StaticFile("/favicon.ico", tz.SiteFile("static/static/favicon.ico"))
	r.StaticFile("/", tz.SiteFile("static/index.html"))
	r.Static("/static", tz.SiteFile("static/static"))
	api := r.Group("/api")
	key1 := api.Group("/keygen")
	{
		key1.POST("charles", keygen.GetCharlesKey)
		key1.GET("appstore", keygen.GetAccounts)
	}
	key2 := r.Group("/rpc")
	{
		key2.GET("/obtainTicket.action", keygen.ActivateIdea)
		key2.GET("/releaseTicket.action", func(*gin.Context) {})
		key2.GET("/prolongTicket.action", func(*gin.Context) {})
	}
	r.NoRoute(func(c *gin.Context) {
		c.File(tz.SiteFile("static/index.html"))
	})
}

func (tz *toolez) SiteName() string {
	return tz.siteName
}

func (tz *toolez) HostNames() []string {
	return tz.config.HostNames
}

func (tz *toolez) GinEngine() *gin.Engine {
	return tz.router
}
