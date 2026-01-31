package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newERPDatabase(dns string, logger gormlogger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(dns), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	// sqlDB.SetMaxIdleConns(cfg.Database.MaxConnections)
	// sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
