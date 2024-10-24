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
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}

	mux.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	mux.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	mux.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe("localhost:8181", mux))

}
