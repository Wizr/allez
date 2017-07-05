package core

import (
	"crypto/tls"
	"net/http"
	"path/filepath"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func NewServer(rootPath string) *EZServer {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(rootPath, "config.yaml")
	c := &Config{RootPath: rootPath}
	return &EZServer{
		siteServers: make(map[string]*EZSiteServer),
		config:      c.ParseFile(configPath),
	}
}

/***************************************************
EZSite the interface of a site
***************************************************/

type EZSite interface {
	Init(config *Config)
	Register(*gin.Engine)
	SiteName() string
	DNSName() []string
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

func (ss *EZSiteServer) DNSName() []string {
	return ss.site.DNSName()
}

/*********************************************************
EZServer a wrapper around gin providing extra features
*********************************************************/

type EZServer struct {
	siteServers map[string]*EZSiteServer
	middleware  []gin.HandlerFunc
	config      *Config

	// certificate
	am      *autocert.Manager
	crtFile *string
	keyFile *string
}

func (s *EZServer) Use(middleware ...gin.HandlerFunc) *EZServer {
	s.middleware = append(s.middleware, middleware...)
	return s
}

func (s *EZServer) SetAutoCert(manager *autocert.Manager) *EZServer {
	s.am = manager
	return s
}

func (s *EZServer) SetCert(crtFile string, keyFile string) *EZServer {
	s.crtFile = &crtFile
	s.keyFile = &keyFile
	return s
}

func (s *EZServer) RegisterSite(sites ...EZSite) *EZServer {
	for _, site := range sites {
		s.siteServers[site.SiteName()] = &EZSiteServer{site: site, server: gin.New()}
	}
	return s
}

func (s *EZServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.config.Mode == "DEV" {
		s.siteServers["toolez"].server.ServeHTTP(w, r)
		return
	}
	if siteServer := s.siteServers[r.Host]; siteServer != nil {
		siteServer.server.ServeHTTP(w, r)
	} else {
		// Handle host names for which no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func (s *EZServer) ListenAndServe() {
	for _, ss := range s.siteServers {
		ss.Init(s.config)
		ss.Register()
		ss.Use(s.middleware...)
	}
	server := endless.NewServer(s.selectAddr(), s)
	if s.config.Mode == "PROD" {
		if s.am != nil {
			hosts := []string{}
			for _, ss := range s.siteServers {
				hosts = append(hosts, ss.DNSName()...)
			}
			s.am.HostPolicy = autocert.HostWhitelist(hosts...)
			server.TLSConfig = &tls.Config{GetCertificate: s.am.GetCertificate}
		}
		server.ListenAndServe()
	} else {
		server.ListenAndServeTLS(*s.crtFile, *s.keyFile)
	}
}

///////////////////
// helper function
func (s *EZServer) selectAddr() string {
	if s.config.Mode == "PROD" {
		return s.config.Prod.Addr
	} else {
		return s.config.Dev.Addr
	}
}

/***************************************************
Custom middleware
***************************************************/
