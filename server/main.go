package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/ComedicChimera/tempest/server/app"
	"github.com/ComedicChimera/tempest/server/controllers"
	"github.com/gorilla/mux"
)

func addHandlers(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/get-nonce", controllers.GetNonce)
}

func main() {
	router := mux.NewRouter()
	addHandlers(router)

	router.Use(app.TransportSecurity)
	router.Use(app.JWTAuthentication)

	srv := &http.Server{
		Addr:         ":443",
		Handler:      router,
		TLSConfig:    app.TLSConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Fatal(srv.ListenAndServeTLS("config/tls.crt", "config/tls.key"))
}
