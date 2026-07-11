package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"taladsync/internal/auth"
	"taladsync/internal/config"
	"taladsync/internal/database"
	"taladsync/internal/domain"
	"taladsync/internal/handler"
	"taladsync/internal/middleware"
	"taladsync/internal/repository"
	"taladsync/internal/usecase"
)

func main() {
	reseed := flag.Bool("reseed", false, "ล้าง DB แล้วใส่ข้อมูล dev สดใหม่ แล้วออก (ไม่ start server)")
	flag.Parse()

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// go run ./cmd/api -reseed  → รีเซ็ตข้อมูลแล้วออก (สำหรับเดโม/dev)
	if *reseed {
		if err := database.Reseed(db); err != nil {
			log.Fatal(err)
		}
		log.Println("reseed complete — exiting")
		return
	}

	if err := database.SeedIfEmpty(db); err != nil {
		log.Fatal(err)
	}
	if err := database.SeedDailyBookingsForToday(db); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenSvc := auth.NewTokenService([]byte(cfg.JWTSecret))
	authUC := usecase.NewAuthUsecase(userRepo, tokenSvc)
	authHandler := handler.NewAuthHandler(authUC)

	bookingRepo := repository.NewBookingRepository(db)
	bookingUC := usecase.NewBookingUsecase(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingUC)

	auditRepo := repository.NewAuditRepository(db)
	vendorUC := usecase.NewVendorUsecase(bookingRepo, auditRepo)
	vendorHandler := handler.NewVendorHandler(vendorUC)
	staffUC := usecase.NewStaffUsecase(bookingRepo, auditRepo)
	staffHandler := handler.NewStaffHandler(staffUC)

	directoryUC := usecase.NewDirectoryUsecase(bookingRepo)
	directoryHandler := handler.NewDirectoryHandler(directoryUC)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Post("/auth/mock-login", authHandler.MockLogin)

	api.Get("/directory", directoryHandler.Search)

	authMW := middleware.Auth(tokenSvc)

	protected := api.Group("/auth", authMW)
	protected.Get("/me", authHandler.Me)


	api.Get("/stalls/available", authMW, middleware.RequireRoles(domain.RoleWalkin), bookingHandler.ListAvailable)
	api.Get("/bookings/mine", authMW, middleware.RequireRoles(domain.RoleWalkin), bookingHandler.GetMine)
	api.Post("/bookings", authMW, middleware.RequireRoles(domain.RoleWalkin), bookingHandler.Book)

	api.Get("/vendor/mine", authMW, middleware.RequireRoles(domain.RoleRegular), vendorHandler.GetMine)
	api.Post("/vendor/leave", authMW, middleware.RequireRoles(domain.RoleRegular), vendorHandler.Leave)

	api.Get("/staff/checklist", authMW, middleware.RequireRoles(domain.RoleStaff), staffHandler.Checklist)
	api.Get("/staff/summary", authMW, middleware.RequireRoles(domain.RoleStaff), staffHandler.Summary)
	api.Post("/bookings/:id/attendance", authMW, middleware.RequireRoles(domain.RoleStaff), staffHandler.MarkAttendance)

	addr := ":" + cfg.AppPort
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
