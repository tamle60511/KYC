package database

import (
	"CQS-KYC/config"
	"CQS-KYC/internal/model"
	"CQS-KYC/logger"
	"context"

	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// 1. Cập nhật Interface: Thêm GetDB() để tương thích code Service cũ
type Database interface {
	DB() *gorm.DB
	GetDB() *gorm.DB
	ERPDB() *gorm.DB
	Close() error
	Ping() error
}

type database struct {
	db    *gorm.DB
	erpDB *gorm.DB
}

func NewDatabase(cfg *config.Config, log *logger.AppLogger) (Database, error) {
	// Config Logger cho GORM
	gormLog := NewLogger(log.Logger).LogMode(gormlogger.Info)

	// 1. Kết nối ERP Database (Dữ liệu nguồn)
	erpDB, err := newERPDatabase(cfg.GetERPDatabaseDSN(), gormLog)
	if err != nil {
		// Có thể chỉ log warning nếu ERP không critical lúc start
		fmt.Printf("⚠️ Warning: Connect to ERP database failed: [%v]\n", err)
	}

	// 2. Kết nối Main Database (Dữ liệu hệ thống EFNET)
	db, err := newDatabase(cfg.GetDSN(), gormLog)
	if err != nil {
		panic(fmt.Sprintf("Connect to main database failed: [%v]", err))
	}

	// 3. === QUAN TRỌNG: AUTO MIGRATE ===
	// Tự động tạo bảng cho Partition, Workflow, Signature, Users...
	err = runMigrations(db)
	if err != nil {
		panic(fmt.Sprintf("Migration failed: [%v]", err))
	}
	fmt.Println("✅ Database Migration completed successfully!")

	return &database{
		db:    db,
		erpDB: erpDB,
	}, nil
}

func MustNewDatabase(cfg *config.Config, logger *logger.AppLogger) Database {
	db, err := NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return db
}

// Hàm khởi tạo kết nối (Giữ nguyên logic của bạn)
func newDatabase(dns string, logger gormlogger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection Pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// --- IMPLEMENT INTERFACE ---

func (d *database) DB() *gorm.DB {
	return d.db
}

// Thêm method này để tương thích với code Service (s.db.GetDB())
func (d *database) GetDB() *gorm.DB {
	return d.db
}

func (d *database) ERPDB() *gorm.DB {
	return d.erpDB
}

func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	// Đóng cả ERP connection nếu cần
	if d.erpDB != nil {
		if sqlERP, err := d.erpDB.DB(); err == nil {
			sqlERP.Close()
		}
	}
	return sqlDB.Close()
}

func (d *database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// --- MIGRATION LOGIC ---
// Hàm này liệt kê TẤT CẢ các bảng cần thiết cho hệ thống
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Request{},
		&model.User{},
		// // 2. Hệ thống Workflow Động (Dynamic Engine)
		&model.WorkflowDefinition{},
		&model.WorkflowStep{},
		&model.WorkflowStepAssignment{},
		&model.WorkflowInstance{},
		&model.WorkflowTask{},
		&model.WorkflowLog{},
		// // 3. Hệ thống Chữ ký điện tử (Signature Trail)
		// &models.DigitalSignature{},
		// &models.SignatureTemplate{},

		// // 4. Hệ thống User & Group (Phân quyền)
		&model.Manager{},
		&model.Position{},
		&model.Department{},

		&model.UserGroup{},
		&model.UserGroupMember{},
	// &models.WorkflowDelegation{},
	)
}
