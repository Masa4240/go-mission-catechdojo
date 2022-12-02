package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("Recover: ", err)
			}
		}()
		fmt.Println("Start Recovery middleware")
		//panic("Panic!")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
