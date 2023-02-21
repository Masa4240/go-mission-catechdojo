package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger, _ := zap.NewProduction()
		defer func(logger *zap.Logger) {
			if err := logger.Sync(); err != nil {
				logger.Info("zap err")
			}
		}(logger)
		// logger.Info("Recovery")
		defer func() {
			err := recover()
			if err != nil {
				logger.Info("Recovered")
				fmt.Println("Recover: ", err)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
