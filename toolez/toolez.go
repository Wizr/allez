package toolez

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/vettu/allez/core"
	"github.com/vettu/allez/toolez/keygen"
)

type tz struct {
	siteName   string
	rootConfig *core.Config
	config     *Config
}

func (tz *tz) File(paths ...string) string {
	paths = append([]string{tz.rootConfig.RootPath}, paths...)
	return filepath.Join(paths...)
}

func (tz *tz) Init(config *core.Config) {
	tz.rootConfig = config
	c := &Config{}
	mapstructure.Decode(config.Site[tz.siteName], c)
	//config.Site[tz.siteName] = c
	tz.config = c
}

func (tz *tz) Register(r *gin.Engine) {
	r.StaticFile("/favicon", tz.File("static/static/favicon.ico"))
	r.StaticFile("/", tz.File("static/index.html"))
	r.Static("/static", tz.File("static/static"))
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
		c.File(tz.File("static/index.html"))
	})
}

func (tz *tz) SiteName() string {
	return tz.siteName
}

func (tz *tz) HostName() []string {
	return tz.config.HostName
}

func New() core.EZSite {
	s := &tz{
		siteName: "toolez",
	}
	return s
}
