package httpserver

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrvin/url-shortener/internal/cache"
	"github.com/mrvin/url-shortener/internal/httpserver/handlers"
	log "github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	"github.com/mrvin/url-shortener/pkg/http/logger"
	"golang.org/x/crypto/bcrypt"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

type ConfHTTPS struct {
	CertFile string
	KeyFile  string
}

type Conf struct {
	Host        string
	Port        string
	IsHTTPS     bool
	HTTPS       ConfHTTPS
	DocFilePath string
}

type Server struct {
	http.Server

	conf *Conf
}

//go:embed static
var staticFiles embed.FS

func New(conf *Conf, st storage.Storage, c cache.Cacher) *Server {
	mux := http.NewServeMux()

	// frontend
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle(http.MethodGet+" /static/", fileServer)
	mux.HandleFunc(http.MethodGet+" /favicon.ico", handlers.GetFavicon)

	// docs
	mux.HandleFunc(http.MethodGet+" /api/openapi.yaml", handlers.NewAPIDocs(conf.DocFilePath))

	// info
	mux.HandleFunc(http.MethodGet+" /api/health", handlers.Health)
	mux.HandleFunc(http.MethodGet+" /api/info", handlers.ErrorHandler("Info", handlers.Info))

	// users
	mux.HandleFunc(http.MethodPost+" /api/users", handlers.ErrorHandler("Registration user", handlers.NewRegistration(st)))
	mux.HandleFunc(http.MethodPost+" /api/users/login", handlers.ErrorHandler("Login user", handlers.NewLogin(st)))

	// urls
	mux.HandleFunc(http.MethodPost+" /api/urls", auth(handlers.ErrorHandler("Save url", handlers.NewSaveURL(st)), st))
	mux.HandleFunc(http.MethodGet+" /api/urls", auth(handlers.ErrorHandler("Get urls", handlers.NewGetURLs(st)), st))
	mux.HandleFunc(http.MethodGet+" /api/urls/check/{alias...}", handlers.ErrorHandler("Check alias", handlers.NewCheckAlias(st)))
	mux.HandleFunc(http.MethodDelete+" /api/urls/{alias...}", auth(handlers.ErrorHandler("Delete url", handlers.NewDeleteURL(st, c)), st))
	mux.HandleFunc(http.MethodGet+" /{alias...}", handlers.ErrorHandler("Redirect", handlers.NewRedirect(st, c)))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustruct
		http.Server{
			Addr:         net.JoinHostPort(conf.Host, conf.Port),
			Handler:      &loggerServer,
			ReadTimeout:  readTimeout * time.Second,
			WriteTimeout: writeTimeout * time.Second,
			IdleTimeout:  idleTimeout * time.Minute,
		},
		conf,
	}
}

func (s *Server) Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)

	go func() {
		defer cancel()
		if s.conf.IsHTTPS {
			if err := s.ListenAndServeTLS(s.conf.HTTPS.CertFile, s.conf.HTTPS.KeyFile); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start https server: " + err.Error())
				return
			}
		} else {
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start http server: " + err.Error())
				return
			}
		}
	}()
	if s.conf.IsHTTPS {
		slog.Info("Start http server: https://" + s.Addr)
	} else {
		slog.Info("Start http server: http://" + s.Addr)
	}

	<-ctx.Done()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Failed to stop http server: " + err.Error())
		return
	}
	slog.Info("Stop http server")
}

type UserGetter interface {
	GetUser(ctx context.Context, name string) (*storage.User, error)
}

func auth(next http.HandlerFunc, getter UserGetter) http.HandlerFunc {
	handler := func(res http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := req.Context()
		user, err := getter.GetUser(ctx, username)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusInternalServerError)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = log.WithUsername(ctx, username)

		next(res, req.WithContext(ctx)) // Pass request to next handler
	}

	return http.HandlerFunc(handler)
}
