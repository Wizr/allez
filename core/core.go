package core

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"golang.org/x/crypto/acme/autocert"
)

func NewServer(rootPath string) *EZServer {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(rootPath, "config.yaml")
	c := &Config{RootPath: rootPath}
	s := &EZServer{
		config: c.ParseFile(configPath),
	}
	if s.config.Mode == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}
	return s
}

/***************************************************
EZSite the interface of a site
***************************************************/

type EZSite interface {
	Init(config *Config)
	Register(*gin.Engine)
	SiteName() string
	HostName() []string
}

type EZSiteServer struct {
	site   EZSite
	server *gin.Engine
}

func (ss *EZSiteServer) Init(c *Config) {
	ss.site.Init(c)
}

func (ss *EZSiteServer) Register() {
	ss.site.Register(ss.server)
}

func (ss *EZSiteServer) Use(middleware ...gin.HandlerFunc) *EZSiteServer {
	ss.server.Use(middleware...)
	return ss
}

func (ss *EZSiteServer) HostName() []string {
	return ss.site.HostName()
}

/*********************************************************
EZServer a wrapper around gin providing extra features
*********************************************************/

type EZServer struct {
	siteServers     []*EZSiteServer
	middleware      []gin.HandlerFunc
	config          *Config
	host2SiteServer map[string]*EZSiteServer

	// certificate
	certManager *autocert.Manager
	crtFile     *string
	keyFile     *string
}

func (s *EZServer) Use(middleware ...gin.HandlerFunc) *EZServer {
	s.middleware = append(s.middleware, middleware...)
	return s
}

func (s *EZServer) SetAutoCert(manager *autocert.Manager) *EZServer {
	s.certManager = manager
	return s
}

func (s *EZServer) SetCert(crtFile string, keyFile string) *EZServer {
	s.crtFile = &crtFile
	s.keyFile = &keyFile
	return s
}

func (s *EZServer) RegisterSite(sites ...EZSite) *EZServer {
	for _, site := range sites {
		s.siteServers = append(s.siteServers, &EZSiteServer{site: site, server: gin.New()})
	}
	return s
}

func (s *EZServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.config.Mode == "DEV" {
		s.siteServers[0].server.ServeHTTP(w, r)
		return
	}
	if siteServer, ok := s.host2SiteServer[r.Host]; ok {
		siteServer.server.ServeHTTP(w, r)
	} else {
		// Handle host names for which no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func (s *EZServer) ListenAndServe() {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:         true,
		SSLRedirect:       true,
		SSLHost:           "lst.nt:8443",
		HostsProxyHeaders: []string{"X-Forwarded-Host"},
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()
	for _, ss := range s.siteServers {
		ss.Init(s.config)
		ss.Use(s.middleware...)
		ss.Use(secureFunc)
		ss.Register()
	}
	serverSSL := &http.Server{Addr: s.selectAddr(true)}
	serverSSL.Handler = s
	if s.config.Mode == "PROD" {
		if s.certManager == nil {
			panic("Please setup cert manager.")
		}
		// setup cert manager
		hostNames := []string{}
		host2SiteServer := make(map[string]*EZSiteServer)
		for _, ss := range s.siteServers {
			for _, host := range ss.HostName() {
				host2SiteServer[host] = ss
			}
			hostNames = append(hostNames, ss.HostName()...)
		}
		s.host2SiteServer = host2SiteServer
		s.certManager.HostPolicy = autocert.HostWhitelist(hostNames...)
		serverSSL.TLSConfig = &tls.Config{GetCertificate: s.certManager.GetCertificate}
		go func() {
			if err := serverSSL.ListenAndServeTLS("", ""); err != nil {
				log.Printf("listen: %s\n", err)
			}
		}()
	} else {
		if _, err := os.Stat(*s.crtFile); os.IsNotExist(err) {
			panic(err)
		}
		if _, err := os.Stat(*s.keyFile); os.IsNotExist(err) {
			panic(err)
		}
		go func() {
			if err := serverSSL.ListenAndServeTLS(*s.crtFile, *s.keyFile); err != nil {
				log.Printf("listen: %s\n", err)
			}
		}()
	}
	go func() {
		server := &http.Server{Addr: s.selectAddr(false)}
		server.Handler = s
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serverSSL.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exited")
}

///////////////////
// helper function
func (s *EZServer) selectAddr(isSSL bool) string {
	if isSSL == false {
		if s.config.Mode == "PROD" {
			return s.config.Prod.Addr
		} else {
			return s.config.Dev.Addr
		}
	} else {
		if s.config.Mode == "PROD" {
			return s.config.Prod.AddrSSL
		} else {
			return s.config.Dev.AddrSSL
		}
	}
}

/***************************************************
Custom middleware
***************************************************/
