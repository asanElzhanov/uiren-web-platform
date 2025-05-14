package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uiren/internal/app/achievements"
	"uiren/internal/app/admin"
	"uiren/internal/app/auth"
	"uiren/internal/app/data"
	"uiren/internal/app/exercises"
	"uiren/internal/app/friendship"
	"uiren/internal/app/lessons"
	"uiren/internal/app/modules"
	"uiren/internal/app/progress"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/database"
	jwt_maker "uiren/internal/infrastracture/jwt"
	yandex_sender "uiren/internal/infrastracture/mail/yandex"
	"uiren/pkg/config"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	//app
	appPortKey  = "app_port"
	appLogLevel = "app_log_level"
	//db postgres
	dbPostgresHostKey = "db_postgres_host"
	dbPostgresPortKey = "db_postgres_port"
	dbPostgresUserKey = "db_postgres_user"
	dbPostgresNameKey = "db_postgres_name"
	//db mongo
	dbMongoName = "db_mongo_name"
	//db redis
	dbRedisAddressKey  = "db_redis_address"
	dbRedisPasswordKey = "db_redis_password"
	dbRedisDBKey       = "db_redis_db"
	dbRedisDataTTLKey  = "db_redis_data_TTL"
	//jwt
	jwtDurationKey       = "jwt_duration"
	refreshTokenDuration = "refresh_token_duration"
	//email
	emailSenderNameKey     = "email_sender_name"
	fromEmailAddressKey    = "from_email_address"
	verificationCodeTTLKey = "verification_code_TTL"
	//data
	xpLeaderboardLimitKey = "xp_leaderboard_limit"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH", // PATCH указан
		AllowHeaders: "Content-Type,Authorization", AllowCredentials: true,
	}))

	err := logger.InitLogger(config.GetValue(appLogLevel).String())
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = godotenv.Load()
	if err != nil {
		logger.Fatal(".env file loading error")
	}

	postgresDB, err := database.GetPostgresDatabase(ctx, database.PostgresConfig{
		PostgresHost: config.GetValue(dbPostgresHostKey).String(),
		PostgresPort: config.GetValue(dbPostgresPortKey).String(),
		PostgresUser: config.GetValue(dbPostgresUserKey).String(),
		PostgresPass: os.Getenv("POSTGRES_DB_PASSWORD"),
		PostgresName: config.GetValue(dbPostgresNameKey).String(),
	})
	if err != nil {
		logger.Fatal("postgres db: ", err)
		return
	}
	logger.Info("connected to database postgres")

	mongoDB, client, err := database.GetMongoDatabase(ctx, database.MongoConfig{
		URI:    os.Getenv("DB_MONGO_URI"),
		DBName: config.GetValue(dbMongoName).String(),
	})
	if err != nil {
		logger.Fatal("mongo db: ", err)
		return
	}
	defer func() {
		if client != nil {
			client.Disconnect(ctx)
		}
	}()
	logger.Info("connected to MongoDB: ", mongoDB.Name())

	redisDB, err := database.GetRedisDatabase(ctx, database.RedisConfig{
		Address:  config.GetValue(dbRedisAddressKey).String(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       config.GetValue(dbRedisDBKey).Int(),
		DataTTL:  config.GetValue(dbRedisDataTTLKey).Duration(),
	})
	if err != nil {
		logger.Fatal("redis db: ", err)
		return
	}
	defer func() {
		if redisDB != nil {
			redisDB.Client.Close()
		}
	}()

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

	achievementRepo := achievements.NewAchievementRepository(postgresDB)
	achievementService := achievements.NewAchievementService(achievementRepo)

	progressReceiverRepo := progress.NewProgressReceiverRepository(postgresDB)
	progressUpdaterRepo := progress.NewProgressUpdaterRepository(postgresDB)
	progressService := progress.NewProgressService(progressReceiverRepo, progressUpdaterRepo, achievementService)

	userRepo := users.NewUserRepository(postgresDB)
	userService := users.NewUserService(userRepo, progressService)

	friendshipRepo := friendship.NewFriendshipRepository(postgresDB)
	friendshipService := friendship.NewFriendshipService(friendshipRepo, userService)

	verifRepo := auth.NewVerificationRepository(postgresDB)
	authService := auth.NewAuthService(userService, jwtMaker, verifRepo)
	authService.SetVerificationCodeTTL(config.GetValue(verificationCodeTTLKey).Duration())
	authService.WithRedisClient(redisDB)
	authService.SetRefreshTokenTTL(config.GetValue(refreshTokenDuration).Duration())

	dataService := data.NewDataService(
		redisDB,
		userService,
		modulesService,
		config.GetValue(dbRedisDataTTLKey).Duration(),
	)
	dataService.WithProgressService(progressService, config.GetValue(xpLeaderboardLimitKey).Int())

	appService := admin.NewApp(app)
	appService.WithUserService(userService)
	appService.WithAuthService(authService)
	appService.WithModulesSerivce(modulesService)
	appService.WithLessonService(lessonService)
	appService.WithExerciseService(exerciseService)
	appService.WithAchievementService(achievementService)
	appService.WithFriendshipService(friendshipService)
	appService.WithDataService(dataService)
	appService.WithProgressService(progressService)
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
