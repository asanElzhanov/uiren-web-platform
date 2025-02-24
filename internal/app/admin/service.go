package admin

import (
	"context"
	"uiren/internal/app/auth"
	"uiren/internal/app/modules"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/middleware"

	"github.com/gofiber/fiber/v2"
)

type modulesService interface {
	GetModule(ctx context.Context, code string) (modules.ModuleDTO, error)
}

type userService interface {
	CreateUser(ctx context.Context, params users.CreateUserDTO) (string, error)
	GetUserForLogin(ctx context.Context, username string) (users.UserDTO, error)
	UpdateUser(ctx context.Context, dto users.UpdateUserDTO) (users.UserDTO, error)
}

type authService interface {
	SignIn(ctx context.Context, params auth.LoginParams) (string, error)
	Register(ctx context.Context, params auth.RegisterParams) (string, error)
	VerifyUser(ctx context.Context, username, code string) error
}

type App struct {
	appFiber       *fiber.App
	userService    userService
	authService    authService
	modulesService modulesService
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

func (app *App) SetHandlers() {
	api := app.appFiber.Group("/api")

	//users
	usersApi := api.Group("/users", middleware.JWTMiddleware())
	usersApi.Post("/", app.createUser)
	usersApi.Get("/:id", app.getUser)
	usersApi.Patch("/:id", app.updateUser)

	//auth
	api.Get("/sign-in", app.signIn)
	api.Get("/register", app.register)
	api.Get("/verify/:username/:code", app.verification)
	//modules
	api.Get("/modules/:code", app.getModuleInfo)
}

//appFiber.Get("/exportData", app.exportData)
//appFiber.Post("/importData", app.importData)
