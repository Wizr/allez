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
	for _, ss := range s.siteServers {
		ss.Init(s.config)
		ss.Register()
		ss.Use(s.middleware...)
	}
	if s.config.Mode == "PROD" {
		server := &http.Server{Addr: s.selectAddr()}
		if s.am != nil {
			hostNames := []string{}
			host2SiteServer := make(map[string]*EZSiteServer)
			for _, ss := range s.siteServers {
				for _, host := range ss.HostName() {
					host2SiteServer[host] = ss
				}
				hostNames = append(hostNames, ss.HostName()...)
			}
			s.host2SiteServer = host2SiteServer
			s.am.HostPolicy = autocert.HostWhitelist(hostNames...)
			server.TLSConfig = &tls.Config{GetCertificate: s.am.GetCertificate}
			server.Handler = s
		}
		server.ListenAndServeTLS("", "")
	} else {
		server := endless.NewServer(s.selectAddr(), s)
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
