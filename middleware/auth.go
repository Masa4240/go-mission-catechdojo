package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/form3tech-oss/jwt-go"
)

type Auth struct {
	Name string
	ID   int64
}

func TokenValidation(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start token validation")
		token := r.Header["X-Token"]

		// fmt.Println(token)
		// fmt.Println(token[0])
		//fmt.Println(r.Header)
		if token[0] == "" {
			fmt.Println("No Token")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dectoken, err := Parse(token[0])

		if err != nil {
			fmt.Println("Fail yo Parse token")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Decode done")

		fmt.Println(dectoken.ID)
		//fmt.Println(dectoken.Name)

		//		h.ServeHTTP(w, r)
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "id", dectoken.ID)))
	}
	return http.HandlerFunc(fn)
}

func Parse(signedString string) (*Auth, error) {
	fmt.Println("start Parse token process")
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//return "", err.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("SIGNINGKEY"), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				fmt.Println(err, "%s is expired", signedString)
				return nil, errors.New("INVALID Token")
			} else {
				fmt.Println(err, "%s is invalid", signedString)
				return nil, errors.New("INVALID Token")
			}
		} else {
			fmt.Println(err, "%s is invalid", signedString)
			return nil, errors.New("INVALID Token")
		}
	}

	if token == nil {
		fmt.Println("not found token in %s:", signedString)
		return nil, errors.New("INVALID Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("not found claims in %s", signedString)
		return nil, errors.New("INVALID Token")
	}
	fmt.Println("Get ID from token")
	fmt.Println(claims)
	id := claims["id"].(float64)
	fmt.Println("ID")

	return &Auth{
		//		Name: name,
		ID: int64(id),
	}, nil
}
