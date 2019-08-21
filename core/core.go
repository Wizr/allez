package core

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/unrolled/secure.v1"
)

func NewServer(rootPath string) *Server {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(rootPath, "config.yaml")
	c := &Config{RootPath: rootPath}
	s := &Server{
		config: c.ParseFile(configPath),
	}
	if s.config.Mode == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.host2SiteInfo = make(map[string]ISiteInfo)
	return s
}

/*********************************************************
Server a wrapper around gin providing extra features
*********************************************************/

type Server struct {
	siteInfos     []ISiteInfo
	middlewares   []gin.HandlerFunc
	config        *Config
	host2SiteInfo map[string]ISiteInfo

	// certificate
	certManager *autocert.Manager
	crtFile     *string
	keyFile     *string
}

func (s *Server) Use(middleware ...gin.HandlerFunc) *Server {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

func (s *Server) SetAutoCert(manager *autocert.Manager) *Server {
	s.certManager = manager
	return s
}

func (s *Server) SetStaticCert(crtFile string, keyFile string) *Server {
	s.crtFile = &crtFile
	s.keyFile = &keyFile
	return s
}

func (s *Server) RegisterSite(sites ...ISiteInfo) *Server {
	for _, site := range sites {
		s.siteInfos = append(s.siteInfos, site)
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if siteServer, ok := s.host2SiteInfo[r.Host]; ok {
		siteServer.GinEngine().ServeHTTP(w, r)
	} else {
		// Handle host names for which no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func (s *Server) ListenAndServe() {
	var allowedHost []string
	for _, si := range s.siteInfos {
		si.Init(s.config)
		for _, host := range si.HostNames() {
			for _, port := range []string{s.getSubConfig().Port, s.getSubConfig().PortSSL} {
				p := ""
				if port != "80" && port != "443" {
					p = ":" + port
				}
				h := host + p
				if _, ok := s.host2SiteInfo[h]; !ok {
					s.host2SiteInfo[h] = si
					allowedHost = append(allowedHost, h)
				}
			}
		}
	}
	s.validateConfig()
	secureFunc := s.getSecureFunc(allowedHost)
	for _, si := range s.siteInfos {
		si.DelayUse(s.middlewares...)
		si.DelayUse(secureFunc)
		si.Use()
		si.RegRouter()
	}
	serverSSL := &http.Server{Addr: s.getListenAddr(true)}
	serverSSL.Handler = s
	if s.config.Mode == "PROD" {
		if s.certManager == nil {
			panic("Please setup cert manager.")
		}
		// setup cert manager
		s.certManager.HostPolicy = autocert.HostWhitelist(allowedHost...)
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
		server := &http.Server{Addr: s.getListenAddr(false)}
		server.Handler = s.certManager.HTTPHandler(s)
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

func (s *Server) getListenAddr(isSSL bool) (addr string) {
	c := s.getSubConfig()
	if isSSL == false {
		addr = c.Port
	} else {
		addr = c.PortSSL
	}
	addr = ":" + addr
	return
}

func (s *Server) getSubConfig() SubConfig {
	return s.config.getSubConfig()
}

func (s *Server) getSecureFunc(allowedHost []string) gin.HandlerFunc {
	sslHostFunc := func() secure.SSLHostFunc {
		cache := make(map[string]string)
		return func(host string) (newHost string) {
			// ignore default port
			if s.getSubConfig().PortSSL == "443" {
				return
			}
			if v, ok := cache[host]; ok {
				newHost = v
			} else {
				slices := strings.Split(host, ":")
				if len(slices) == 2 {
					newHost = slices[0] + ":" + s.getSubConfig().PortSSL
					cache[host] = newHost
				}
			}
			return
		}
	}()
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:   true,
		SSLRedirect: true,
		SSLHostFunc: &sslHostFunc,
		//AllowedHosts: allowedHost,
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
	return secureFunc
}

func (s *Server) validateConfig() {
	s.config.Validate()
	for _, rootConfig := range s.config.Site {
		if siteConfig, ok := rootConfig.(iValidate); ok {
			siteConfig.Validate()
		}
	}
}
