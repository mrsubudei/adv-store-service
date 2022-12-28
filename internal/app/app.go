package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrsubudei/adv-store-service/internal/config"
	v1 "github.com/mrsubudei/adv-store-service/internal/controller/http/v1"
	"github.com/mrsubudei/adv-store-service/internal/repository/sqlite"
	"github.com/mrsubudei/adv-store-service/internal/service"

	"github.com/mrsubudei/adv-store-service/pkg/httpserver"
	"github.com/mrsubudei/adv-store-service/pkg/logger"
	"github.com/mrsubudei/adv-store-service/pkg/sqlite3"
)

func Run(cfg config.Config) {
	// Logger
	l := logger.New()

	// Sqlite
	sq, err := sqlite3.New("database/adverts.db")
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - sqlite3.New: %w", err))
		return
	}
	defer sq.Close()

	// Repository
	repo := sqlite.NewAdvertsRepo(sq)
	err = sqlite.CreateDB(sq)
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - NewRepositories: %w", err))
		return
	}

	// Service
	service := service.NewAdvertService(repo)

	// Http
	handler := v1.NewHandler(service, cfg, l)
	server := httpserver.NewServer(handler)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	fmt.Printf("Server started at http://%s%s\n", cfg.Server.Host, cfg.Server.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	err = server.Shutdown()
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
