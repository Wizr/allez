package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/vettu/allez/core"
	"github.com/vettu/allez/libs/middleware"
	"github.com/vettu/allez/toolez"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

var rootPath *string

func init() {
	p := filepath.Dir(os.Args[0])
	rootPath = flag.String("root", p, "Root path, default is the executable path.")
	flag.Parse()
}

func main() {
	core.NewServer(*rootPath).
		Use(gin.Logger(), gin.Recovery(), middleware.Redis()).
		RegisterSite(toolez.New()).
		SetAutoCert(&autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("./certcache"),
		}).
		SetStaticCert("./certcache/server.crt", "./certcache/server.key").
		ListenAndServe()
}
