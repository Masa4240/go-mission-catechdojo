package router

import (
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/controller"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service"
	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

func NewRouter(userDB *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)

	svc3 := usermodel.NewUserModel(userDB)
	svc2 := userservice.NewUserService(svc3)
	userController := controller.NewUserController(svc2)
	r.Route("/user1", func(r chi.Router) {
		r.Post("/create", userController.CreateUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.TokenValidation)
			r.Get("/get", userController.GetUser)
			r.Put("/update", userController.UpdateUser)
		})
	})

	gsvc3 := gachamodel.NewGachaModel(userDB)
	gsvc2 := gachaservice.NewGachaService(gsvc3)
	gachaController := controller.NewGachaController(gsvc2)
	r.Route("/gacha1", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Post("/draw", gachaController.Gacha)
	})

	r.Route("/character1", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Get("/list", gachaController.GetCharacterList)
	})
	// For admin Usage
	r.Post("/char/add", gachaController.AddCharacter)

	return r
}
