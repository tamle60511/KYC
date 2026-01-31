package app

import (
	"CQS-KYC/config"
	"CQS-KYC/database"
	"CQS-KYC/internal/handler"
	"CQS-KYC/internal/repository"
	"CQS-KYC/internal/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	flogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// Interface này để đảm bảo mọi Handler đều có hàm SetupRoutes
type RouteHandler interface {
	SetupRoutes(router fiber.Router)
}

type App struct {
	config      *config.Config
	fiber       *fiber.App
	database    database.Database
	handlers    []handler.BaseHandler // Danh sách REST Handlers
	soapHandler *handler.SOAPHandler  // Handler riêng cho ERP (SOAP)
}

func New(cfg *config.Config, db database.Database) *App {
	app := &App{
		config:   cfg,
		database: db,
	}

	app.fiber = fiber.New(fiber.Config{
		AppName:      cfg.Server.Name,
		ErrorHandler: errorHandler,
	})

	// --- Middlewares ---
	app.fiber.Use(recover.New())
	app.fiber.Use(flogger.New())
	app.fiber.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Dev only, Pro nên siết lại
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))

	// =========================================================================
	// DEPENDENCY INJECTION (KHỞI TẠO CÁC LỚP)
	// =========================================================================

	// 1. Utils & Core DB
	gormDB := app.database.DB()
	sigHelper := repository.NewSignatureHelper(&app.config.SignatureKey)

	// 2. Repositories
	userRepo := repository.NewUserRepo(gormDB)
	departmentRepo := repository.NewDepartmentRepo(gormDB)
	managerRepo := repository.NewManagerRepo(gormDB)
	positionRepo := repository.NewPositionRepo(gormDB)
	factoryRepo := repository.NewFactoryRepo(gormDB)
	groupRepo := repository.NewGroupRepo(gormDB)

	wfDefRepo := repository.NewWorkflowRepo(gormDB)

	// Engine cần: DB, GroupRepo (để tìm nhóm), SignatureHelper (để ký)
	instanceRepo := repository.NewWorkflowEngine(gormDB, groupRepo, *sigHelper)

	// 3. Services
	// Service quản lý định nghĩa quy trình (CRUD Workflow)
	wfDefService := service.NewWorkflowService(wfDefRepo)
	departmentService := service.NewDepartmentSerivce(departmentRepo)
	managerService := service.NewManagerService(managerRepo)
	positionService := service.NewPositionService(positionRepo)
	// Service quản lý chạy luồng (Engine)
	instanceService := service.NewInstanceService(instanceRepo, gormDB)
	groupService := service.NewGroupService(groupRepo)
	userService := service.NewUserService(userRepo)
	factoryService := service.NewFactoryService(factoryRepo)
	// Service ERP (Cầu nối)
	erpService := service.NewERPService(app.database, cfg, userRepo, wfDefService, instanceService)

	// 4. Handlers
	// Handler cho REST API (Frontend gọi)
	wfDefHandler := handler.NewWorkflowHandler(wfDefService)
	instanceHandler := handler.NewInstanceHandler(instanceService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	managerHandler := handler.NewManagerHandler(managerService)
	positionHandler := handler.NewPositionHandler(positionService)
	userHandler := handler.NewUserHandler(userService)
	groupHandler := handler.NewGroupHandler(groupService)
	factoryHandler := handler.NewFactoryHandler(factoryService)
	// Handler cho SOAP API (ERP gọi)
	soapHandler := handler.NewSOAPHandler(erpService)

	app.handlers = []handler.BaseHandler{
		userHandler,
		groupHandler,
		factoryHandler,
		wfDefHandler,
		instanceHandler,
		departmentHandler,
		managerHandler,
		positionHandler,
	}
	app.soapHandler = soapHandler
	return app
}

func (a *App) SetupRoutes() {
	// 1. Health Check
	a.fiber.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"name":   a.config.Server.Name,
			"env":    a.config.Server.ENV,
		})
	})

	// =========================================================================
	// 2. SOAP ROUTE CHO ERP
	// =========================================================================
	if a.soapHandler != nil {
		soapGroup := a.fiber.Group("/EFNETService")
		// Endpoint nhận XML từ ERP
		soapGroup.Post("/EFERPService.asmx", a.soapHandler.HandleRequest)
		// Endpoint trả về document WSDL (để ERP biết cách gọi)
		soapGroup.Get("/EFERPService.asmx", a.soapHandler.HandleWSDL)
	}

	// =========================================================================
	// 3. REST API (Cho Frontend React/Mobile)
	// =========================================================================
	api := a.fiber.Group("/api")

	// Tự động setup route cho tất cả handler trong list
	for _, h := range a.handlers {
		// Lưu ý: Các handler cần implement method: SetupRoutes(router fiber.Router)
		h.SetupRoutes(api)
	}

	// 4. Fallback 404
	a.fiber.Use(func(c fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Not Found",
			"path":    c.Path(),
		})
	})
}

func (a *App) Start() {
	// Graceful Shutdown Setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", a.config.Server.Port)
		if err := a.fiber.Listen(addr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("🚀 Server started on port %s", a.config.Server.Port)
	log.Printf("📡 SOAP Endpoint: http://localhost:%s/EFNETService/EFERPService.asmx", a.config.Server.Port)
	log.Printf("🔌 REST API: http://localhost:%s/api", a.config.Server.Port)

	// Block main thread until signal received
	<-sigChan
	log.Println("Shutting down server...")

	if err := a.database.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	if err := a.fiber.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}
