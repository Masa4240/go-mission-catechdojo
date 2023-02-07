package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

type Auth struct {
	Name string
	ID   int64
}

func TokenValidation(h http.Handler) http.Handler {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start token validation", zap.Time("now", time.Now()))

	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["X-Token"]

		if token == nil {
			logger.Info("No token", zap.Time("now", time.Now()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if token[0] == "" {
			logger.Info("No token", zap.Time("now", time.Now()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dectoken, err := parse(token[0])
		if err != nil {
			logger.Info("Fail yo Parse token", zap.Time("now", time.Now()), zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Info("Decode Done", zap.Time("now", time.Now()),
			zap.String("Name", dectoken.Name), zap.Int64("ID", dectoken.ID))

		id := "id"
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), id, dectoken.ID)))
	}
	return http.HandlerFunc(fn)
}

func parse(signedString string) (*Auth, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Info("Token Parse Failure", zap.Time("now", time.Now()))
			// return "", err.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("SIGNINGKEY"), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				logger.Info("token is expired", zap.Time("now", time.Now()), zap.Error(err))
			} else {
				logger.Info("token is invalid", zap.Time("now", time.Now()), zap.Error(err))
				return nil, errors.New("INVALID Token")
			}
		} else {
			logger.Info("token is expired", zap.Time("now", time.Now()), zap.Error(err))
			return nil, errors.New("INVALID Token")
		}
	}

	if token == nil {
		logger.Info("not found token :", zap.Time("now", time.Now()), zap.Error(err))
		return nil, errors.New("INVALID Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Info("not found claims in token", zap.Time("now", time.Now()), zap.Error(err))
		return nil, errors.New("INVALID Token")
	}
	id := claims["id"].(float64)

	return &Auth{
		//		Name: name,
		ID: int64(id),
	}, nil
}
