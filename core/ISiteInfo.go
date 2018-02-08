package core

import "github.com/gin-gonic/gin"

type ISiteInfo interface {
	Init(*Config)
	DelayUse(...gin.HandlerFunc) ISiteInfo
	Use()
	RegRouter()
	SiteName() string
	HostNames() []string
	GinEngine() *gin.Engine
}
