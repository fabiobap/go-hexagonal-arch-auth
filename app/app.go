package app

import (
	"log"
	"net/http"

	"github.com/go-hexagonal-arch-auth/domain"
	"github.com/go-hexagonal-arch-auth/service"
	"github.com/gorilla/mux"
)

func Start() {
	mux := mux.NewRouter()

	dbClient := getDBClient()

	authRepository := domain.NewAuthRepository(dbClient)
	ah := AuthHandler{service.NewLoginService(authRepository)}

	mux.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe("localhost:8181", mux))

}
