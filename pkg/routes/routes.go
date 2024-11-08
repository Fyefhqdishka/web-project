package routes

import (
	"database/sql"
	"github.com/Fyefhqdishka/web-project/internal/auth"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func RegisterRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	Auth(r, db, logger)
	HomePage(r)
	TursPage(r)
	StaticFiles(r)
	SignUp(r)
	SignIn(r)
}

func Auth(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	authRepo := auth.NewRepository(db, logger)
	authCtrl := auth.NewControllerAuth(authRepo, logger)

	r.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	r.HandleFunc("/api/login", authCtrl.Login).Methods("POST")

}

func HomePage(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/ui/index.html")
	})
}

func SignIn(r *mux.Router) {
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/ui/login.html")
	})
}

func SignUp(r *mux.Router) {
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/ui/auth.html")
	})
}

func TursPage(r *mux.Router) {
	r.HandleFunc("/turs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/ui/population-turs.html")
	})
}

func StaticFiles(r *mux.Router) {
	fs := http.FileServer(http.Dir("./internal/ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}
