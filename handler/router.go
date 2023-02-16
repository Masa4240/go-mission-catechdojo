package router

import (
	"net/http"

	"github.com/jinzhu/gorm"

	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
	userhandler "github.com/Masa4240/go-mission-catechdojo/handler/user"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"

	gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	gachahandler "github.com/Masa4240/go-mission-catechdojo/handler/gacha"
	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	"github.com/go-chi/chi"
)

func NewRouter(userDB *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)

	userModel := usermodel.NewUserModel(userDB)
	// userService := userservice.NewUserService(userModel)
	userCtrl := usercontroller.NewUserController(userModel)
	userHandler := userhandler.NewUserHandler(userCtrl)

	r.Route("/user", func(r chi.Router) {
		r.Post("/create", userHandler.CreateUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.TokenValidation)
			r.Get("/get", userHandler.GetUser)
			r.Put("/update", userHandler.UpdateUser)
		})
	})

	gachaModel := gachamodel.NewGachaModel(userDB)
	// userService := userservice.NewUserService(userModel)
	gachaCtrl := gachacontroller.NewGachaController(gachaModel)
	gachaHandler := gachahandler.NewGachaHandler(gachaCtrl)

	r.Route("/gacha", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Post("/draw", gachaHandler.Gacha)
	})

	r.Route("/character", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Get("/list", gachaHandler.GetCharacterList)
	})
	// For admin Usage
	r.Post("/char/add", gachaHandler.AddCharacter)

	return r
}
