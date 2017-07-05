package main

import (
	"flag"
	"os"
	"path/filepath"

	"git.oschina.net/nt6/allez/core"
	"git.oschina.net/nt6/allez/toolez"

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
		Use(gin.Logger(), gin.Recovery()).
		RegisterSite(toolez.New()).
		SetAutoCert(&autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("./certcache"),
		}).
		SetCert("./certcache/server.crt", "./certcache/server.key").
		ListenAndServe()
}
