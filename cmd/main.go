package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uiren/internal/app/admin"
	"uiren/internal/app/auth"
	"uiren/internal/app/exercises"
	"uiren/internal/app/lessons"
	"uiren/internal/app/modules"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/database"
	jwt_maker "uiren/internal/infrastracture/jwt"
	yandex_sender "uiren/internal/infrastracture/mail/yandex"
	"uiren/pkg/config"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	//app
	appPortKey = "app_port"
	//db postgres
	dbPostgresHostKey = "db_postgres_host"
	dbPostgresPortKey = "db_postgres_port"
	dbPostgresUserKey = "db_postgres_user"
	dbPostgresNameKey = "db_postgres_name"
	//db mongo
	dbMongoName = "db_mongo_name"
	//jwt
	jwtDurationKey = "jwt_duration"
	//email
	emailSenderNameKey     = "email_sender_name"
	fromEmailAddressKey    = "from_email_address"
	verificationCodeTTLKey = "verification_code_TTL"
)

func main() {
	app := fiber.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file loading error")
	}

	postgresDB, err := database.GetPostgresDatabase(ctx, database.PostgresConfig{
		PostgresHost: config.GetValue(dbPostgresHostKey).String(),
		PostgresPort: config.GetValue(dbPostgresPortKey).String(),
		PostgresUser: config.GetValue(dbPostgresUserKey).String(),
		PostgresPass: os.Getenv("POSTGRES_DB_PASSWORD"),
		PostgresName: config.GetValue(dbPostgresNameKey).String(),
	})
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Info("connected to database postgres")

	mongoDB, client, err := database.GetMongoDatabase(ctx, database.MongoConfig{
		URI:    os.Getenv("DB_MONGO_URI"),
		DBName: config.GetValue(dbMongoName).String(),
	})
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer func() {
		if client != nil {
			client.Disconnect(ctx)
		}
	}()
	logger.Info("connected to MongoDB: ", mongoDB.Name())

	yandex_sender.Init(
		config.GetValue(emailSenderNameKey).String(),
		config.GetValue(fromEmailAddressKey).String(),
		os.Getenv("YANDEX_EMAIL_PASSWORD"),
	)

	jwtMaker := jwt_maker.NewJWTMaker(config.GetValue(jwtDurationKey).Duration())

	exerciseRepo := exercises.NewExercisesRepository(mongoDB)
	exerciseService := exercises.NewExerciseService(exerciseRepo)

	lessonRepo := lessons.NewLessonRepository(mongoDB)
	lessonService := lessons.NewLessonsService(lessonRepo, exerciseService)

	moduleRepo := modules.NewModulesRepository(mongoDB)
	modulesService := modules.NewModulesService(moduleRepo, lessonService)

	userRepo := users.NewUserRepository(postgresDB)
	userService := users.NewUserService(userRepo)

	verifRepo := auth.NewVerificationRepository(postgresDB)
	authService := auth.NewAuthService(userService, jwtMaker, verifRepo)
	authService.SetVerificationCodeTTL(config.GetValue(verificationCodeTTLKey).Duration())

	appService := admin.NewApp(app)
	appService.WithUserService(userService)
	appService.WithAuthService(authService)
	appService.WithModulesSerivce(modulesService)
	appService.SetHandlers()

	port := config.GetValue(appPortKey).String()
	serverErrChan := make(chan error)

	go func() {
		logger.Info("application starting on port ", port)
		serverErrChan <- app.Listen(fmt.Sprintf(":%s", port))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.Info("shutting down application...")
	case err := <-serverErrChan:
		if err != nil {
			logger.Error("server error: ", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(shutdownCtx)
	logger.Info("application closed")
}
