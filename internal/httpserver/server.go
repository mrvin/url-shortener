package httpserver

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrvin/tasks-go/url-shortener/internal/httpserver/handlers"
	log "github.com/mrvin/tasks-go/url-shortener/internal/logger"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/logger"
	"golang.org/x/crypto/bcrypt"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

//nolint:tagliatelle
type ConfHTTPS struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

//nolint:tagliatelle
type Conf struct {
	Addr    string    `yaml:"addr"`
	IsHTTPS bool      `yaml:"is_https"`
	HTTPS   ConfHTTPS `yaml:"https"`
}

type Server struct {
	http.Server

	conf *Conf
}

//go:embed static
var staticFiles embed.FS

func New(conf *Conf, st storage.Storage) *Server {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle(http.MethodGet+" /static/", fileServer)
	mux.HandleFunc(http.MethodGet+" /favicon.ico", handlers.GetFavicon)

	mux.HandleFunc(http.MethodGet+" /api/health", handlers.Health)

	mux.HandleFunc(http.MethodPost+" /api/users", handlers.NewRegistration(st))
	mux.HandleFunc(http.MethodPost+" /api/login", handlers.NewLogin(st))

	mux.HandleFunc(http.MethodPost+" /api/data/shorten", auth(handlers.NewSaveURL(st), st))
	mux.HandleFunc(http.MethodGet+" /api/urls", auth(handlers.NewGetURLs(st), st))
	mux.HandleFunc(http.MethodDelete+" /api/{alias...}", auth(handlers.NewDeleteURL(st), st))

	mux.HandleFunc(http.MethodGet+" /api/check/{alias...}", handlers.NewCheckAlias(st))

	mux.HandleFunc(http.MethodGet+" /{alias...}", handlers.NewRedirect(st))

	loggerServer := logger.Logger{Inner: mux}

	return &Server{
		//nolint:exhaustruct
		http.Server{
			Addr:         conf.Addr,
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
	defer cancel()

	go func() {
		if s.conf.IsHTTPS {
			if err := s.ListenAndServeTLS(s.conf.HTTPS.CertFile, s.conf.HTTPS.KeyFile); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start https server: " + err.Error())
				defer cancel()
				return
			}
		} else {
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start http server: " + err.Error())
				defer cancel()
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
