package router

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
	userhandler "github.com/Masa4240/go-mission-catechdojo/handler/user"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	rankmodel "github.com/Masa4240/go-mission-catechdojo/model/rankratio"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"

	uccontroller "github.com/Masa4240/go-mission-catechdojo/controller/usercharacter"
	uchandler "github.com/Masa4240/go-mission-catechdojo/handler/usercharacter"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	ucservice "github.com/Masa4240/go-mission-catechdojo/service/usercharacter"

	// gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	gachahandler "github.com/Masa4240/go-mission-catechdojo/handler/gacha"
	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"

	"github.com/go-chi/chi"
)

func NewRouter(userDB *gorm.DB, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)

	userModel := usermodel.NewUserModel(userDB, logger)
	userService := userservice.NewUserService(userModel, logger)
	userCtrl := usercontroller.NewUserController(userService, logger)
	userHandler := userhandler.NewUserHandler(userCtrl, logger)

	r.Route("/user", func(r chi.Router) {
		r.Post("/create", userHandler.CreateUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.TokenValidation)
			r.Get("/get", userHandler.GetUser)
			r.Put("/update", userHandler.UpdateUser)
		})
	})

	ucModel := ucmodel.NewUcModel(userDB, logger)
	ucService := ucservice.NewUcService(ucModel, logger)
	ucCtrl := uccontroller.NewUcController(ucService, logger)
	ucHandler := uchandler.NewUcHandler(ucCtrl, logger)

	r.Route("/character", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Get("/list", ucHandler.GetCharacterList)
	})
	cModel := charactermodel.NewCharacterModel(userDB, logger)
	rModel := rankmodel.NewRankModel(userDB, logger)

	gachaService := gachaservice.NewGachaService(cModel, ucModel, userModel, rModel, logger)
	gachaCtrl := gachacontroller.NewGachaController(gachaService, logger)
	gachaHandler := gachahandler.NewGachaHandler(gachaCtrl, logger)

	r.Route("/gacha", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Post("/draw", gachaHandler.Gacha)
	})
	// // For admin Usage
	// r.Post("/char/add", gachaHandler.AddCharacter)

	return r
}
