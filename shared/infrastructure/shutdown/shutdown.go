package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"gorm.io/gorm"
)

type Handler struct {
	app         *fiber.App
	db          *gorm.DB
	redisClient *redis.RedisClient
	timeout     time.Duration
}

func NewHandler(app *fiber.App, db *gorm.DB, redisClient *redis.RedisClient) *Handler {
	return &Handler{
		app:         app,
		db:          db,
		redisClient: redisClient,
		timeout:     15 * time.Second,
	}
}

func (h *Handler) WithTimeout(duration time.Duration) *Handler {
	h.timeout = duration
	return h
}

func (h *Handler) Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("[Shutdown] Signal received. Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	start := time.Now()

	// Shutdown Fiber HTTP Server
	if err := h.app.ShutdownWithContext(ctx); err != nil {
		log.Printf("[Shutdown] Error shutting down HTTP server: %v", err)
	}

	// Shutdown Redis
	h.shutdownRedis()

	// Shutdown Database
	h.shutdownDB(ctx)

	duration := time.Since(start)
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("[Shutdown] Timeout exceeded. Forced shutdown after %s", duration)
	} else {
		log.Printf("[Shutdown] Completed in %s", duration)
	}
}

func (h *Handler) shutdownRedis() {
	if err := h.redisClient.CacheClient.Shutdown(); err != nil {
		log.Printf("[Shutdown] Redis Cache error: %v", err)
	} else {
		log.Println("[Shutdown] Redis Cache closed")
	}

	if err := h.redisClient.LockClient.Shutdown(); err != nil {
		log.Printf("[Shutdown] Redis Lock error: %v", err)
	} else {
		log.Println("[Shutdown] Redis Lock closed")
	}
}

func (h *Handler) shutdownDB(ctx context.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		log.Printf("[Shutdown] Failed to get DB instance: %v", err)
		return
	}

	sqlDB.SetMaxOpenConns(0)
	log.Println("[Shutdown] DB: Stopped new connections")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Shutdown] DB timeout. Forcing close")
			_ = sqlDB.Close()
			return
		case <-ticker.C:
			if sqlDB.Stats().InUse == 0 {
				if err := sqlDB.Close(); err != nil {
					log.Printf("[Shutdown] DB close error: %v", err)
				} else {
					log.Println("[Shutdown] DB closed cleanly")
				}
				return
			} else {
				log.Printf("[Shutdown] DB: Waiting on %d connections", sqlDB.Stats().InUse)
			}
		}
	}
}
