package handler

import (
	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"net/http"

	gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	"github.com/jinzhu/gorm"

	userhandler "github.com/Masa4240/go-mission-catechdojo/handler/user"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"
	"github.com/go-chi/chi"
)

func NewRouter(userDB *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)

	userModel := usermodel.NewUserModel(userDB)
	userService := userservice.NewUserService(userModel)
	userCtrl := usercontroller.NewUserController(userService)

	userHandler := userhandler.NewUserHandler(userCtrl)
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", userHandler.CreateUser)
	})

	gsvc3 := gachamodel.NewGachaModel(userDB)
	gsvc2 := gachaservice.NewGachaService(gsvc3)
	gachaController := gachacontroller.NewGachaController(gsvc2)
	r.Route("/gacha", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Post("/draw", gachaController.Gacha)
	})

	r.Route("/character", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Get("/list", gachaController.GetCharacterList)
	})
	// For admin Usage
	r.Post("/char/add", gachaController.AddCharacter)

	return r
}
