package admin

import (
	"context"
	"uiren/internal/app/achievements"
	"uiren/internal/app/auth"
	"uiren/internal/app/exercises"
	"uiren/internal/app/friendship"
	"uiren/internal/app/lessons"
	"uiren/internal/app/modules"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fiberInternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": ErrInternalServerError})
}

func fiberOK(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "OK"})
}

type modulesService interface {
	GetModule(ctx context.Context, code string) (modules.ModuleWithLessons, error)
	CreateModule(ctx context.Context, dto modules.CreateModuleDTO) (primitive.ObjectID, error)
	DeleteModule(ctx context.Context, code string) error
	UpdateModule(ctx context.Context, code string, dto modules.UpdateModuleDTO) error
	AddLessonToList(ctx context.Context, code, lessonCode string) error
	DeleteLessonFromList(ctx context.Context, code, lessonCode string) error
}

type lessonService interface {
	GetLesson(ctx context.Context, code string) (lessons.LessonDTO, error)
	CreateLesson(ctx context.Context, dto lessons.CreateLessonDTO) (primitive.ObjectID, error)
	UpdateLesson(ctx context.Context, code string, dto lessons.UpdateLessonDTO) error
	DeleteLesson(ctx context.Context, code string) error
	AddExerciseToList(ctx context.Context, code, exerciseCode string) error
	DeleteExerciseFromList(ctx context.Context, code, exerciseCode string) error
}

type exerciseService interface {
	GetExercise(ctx context.Context, code string) (exercises.Exercise, error)
	CreateExercise(ctx context.Context, dto exercises.CreateExerciseDTO) (primitive.ObjectID, error)
	UpdateExercise(ctx context.Context, code string, dto exercises.UpdateExerciseDTO) error
	DeleteExercise(ctx context.Context, code string) error
}

type achievementService interface {
	CreateAchievement(ctx context.Context, name string) (achievements.AchievementDTO, error)
	GetAchievement(ctx context.Context, id int) (achievements.AchievementDTO, error)
	UpdateAchievement(ctx context.Context, dto achievements.UpdateAchievementDTO) (string, error)
	DeleteAchievement(ctx context.Context, id int) error

	AddAchievementLevel(ctx context.Context, dto achievements.AddAchievementLevelDTO) error
	DeleteAchievementLevel(ctx context.Context, dto achievements.DeleteAchievementLevelDTO) error
}

type userService interface {
	CreateUser(ctx context.Context, params users.CreateUserDTO) (string, error)
	GetUserForLogin(ctx context.Context, username string) (users.UserDTO, error)
	UpdateUser(ctx context.Context, dto users.UpdateUserDTO) (users.UserDTO, error)
}

type authService interface {
	SignIn(ctx context.Context, params auth.LoginParams) (string, string, error)
	Register(ctx context.Context, params auth.RegisterParams) (string, error)
	VerifyUser(ctx context.Context, username, code string) error
	RefreshToken(ctx context.Context, token string) (string, string, error)
}

type friendshipService interface {
	SendFriendRequest(ctx context.Context, friendshipRequest friendship.FriendshipRequestDTO) (friendship.Friendship, error)
	HandleFriendRequest(ctx context.Context, friendshipRequest friendship.FriendshipRequestDTO) (friendship.Friendship, error)
	GetFriendList(ctx context.Context, username string) (friendship.FriendList, error)
	GetRequestList(ctx context.Context, username string) (friendship.FriendList, error)
}

type App struct {
	appFiber           *fiber.App
	userService        userService
	authService        authService
	modulesService     modulesService
	lessonService      lessonService
	exerciseService    exerciseService
	achievementService achievementService
	friendshipService  friendshipService
}

func NewApp(appFiber *fiber.App) *App {
	return &App{
		appFiber: appFiber,
	}
}

func (app *App) WithUserService(userService userService) {
	app.userService = userService
}

func (app *App) WithAuthService(authService authService) {
	app.authService = authService
}

func (app *App) WithModulesSerivce(modulesService modulesService) {
	app.modulesService = modulesService
}

func (app *App) WithLessonService(lessonService lessonService) {
	app.lessonService = lessonService
}

func (app *App) WithExerciseService(exerciseService exerciseService) {
	app.exerciseService = exerciseService
}

func (app *App) WithAchievementService(achievementService achievementService) {
	app.achievementService = achievementService
}

func (app *App) WithFriendshipService(friendshipService friendshipService) {
	app.friendshipService = friendshipService
}

func (app *App) SetHandlers() {
	api := app.appFiber.Group("/api")

	//auth
	api.Get("/sign-in", app.signIn)
	api.Get("/register", app.register)
	api.Get("/verify/:username/:code", app.verification)
	api.Get("/refresh-token", app.refreshToken)
	//users
	usersApi := api.Group("/users", middleware.JWTMiddleware(), middleware.AdminMiddleware())
	usersApi.Get("/:id", app.getUser)
	usersApi.Post("/", app.createUser)
	usersApi.Patch("/:id", app.updateUser)
	//modules
	modulesApi := api.Group("/modules", middleware.JWTMiddleware(), middleware.AdminMiddleware())
	modulesApi.Get("/:code", app.getModule)
	modulesApi.Post("/", app.createModule)
	modulesApi.Delete("/:code", app.deleteModule)
	modulesApi.Patch("/:code", app.updateModule)
	modulesApi.Post("/:code/lessons-list/:lessonCode", app.addLessonToList)
	modulesApi.Delete("/:code/lessons-list/:lessonCode", app.deleteLessonFromList)
	//lessons
	lessonApi := api.Group("/lessons", middleware.JWTMiddleware(), middleware.AdminMiddleware())
	lessonApi.Get("/:code", app.getLesson)
	lessonApi.Post("/", app.createLesson)
	lessonApi.Patch("/:code", app.updateLesson)
	lessonApi.Delete("/:code", app.deleteLesson)
	lessonApi.Post(":code/exercises-list/:exerciseCode", app.addExerciseToList)
	lessonApi.Delete(":code/exercises-list/:exerciseCode", app.deleteExerciseFromList)
	//exercises
	exerciseApi := api.Group("/exercises", middleware.JWTMiddleware(), middleware.AdminMiddleware())
	exerciseApi.Get("/:code", app.getExercise)
	exerciseApi.Post("/", app.createExercise)
	exerciseApi.Patch("/:code", app.updateExercise)
	exerciseApi.Delete("/:code", app.deleteExercise)
	//achievements
	achievementsApi := api.Group("/achievements", middleware.JWTMiddleware(), middleware.AdminMiddleware())
	achievementsApi.Post("/", app.createAchievement)
	achievementsApi.Patch("/", app.updateAchievement)
	achievementsApi.Delete("/", app.deleteAchievement)
	achievementsApi.Get("/:id", app.getAchievement)
	achievementsApi.Post("/levels", app.addAchievementLevel)
	achievementsApi.Delete("/levels", app.deleteAchievementLevel)
	//friends
	friendsApi := api.Group("/friends", middleware.JWTMiddleware())
	friendsApi.Post("/send-request", app.sendFriendRequest)
	friendsApi.Post("/handle-request", app.handleFriendRequest)
	friendsApi.Get("/friend-list", app.getFriendList)
	friendsApi.Get("/request-list", app.getRequestList)
}
