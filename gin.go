package log

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	gormzap "github.com/hysios/gorm-zap"
	"github.com/jinzhu/gorm"
)

func LogMiddleware() gin.HandlerFunc {
	return ginzap.Ginzap(logger.Logger, time.RFC3339, true)
}

func SetupGin(r *gin.Engine) {
	r.Use(LogMiddleware())
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))
}

func SetupGorm(db *gorm.DB) {
	db.SetLogger(gormzap.New(logger.Logger))
}
