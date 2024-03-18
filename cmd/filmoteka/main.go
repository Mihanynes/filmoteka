package main

import (
	"database/sql"
	"filmoteka/internal/actor"
	"filmoteka/internal/auth"
	"filmoteka/internal/film"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

const (
	host     = "dbPostgres"
	port     = 5432
	user     = "postgres"
	password = "111111"
	dbname   = "postgres"
)

//@title Filmoteka API
//@version 1.0
//@description This is a Filmoteka server.

// @host localhost:8080
// @basePath /
func main() {
	dsn := "host=dbPostgres port=5432 user=postgres dbname=postgres password=111111 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	actorRepo := actor.NewActorRepository(db)
	filmRepo := film.NewFilmRepository(actorRepo, db)

	a := actor.ActorHandler{
		ActorRepo: actorRepo,
	}
	f := film.FilmHandler{
		FilmRepo: filmRepo,
	}

	sm := auth.NewSessionsDB(db)

	u := &auth.UserHandler{
		DB:       db,
		Sessions: sm,
	}

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Отдельный маршрут для JSON-документации Swagger без аутентификации
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/docs/swagger.json")
	})

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/actor/delete", a.DeleteActor)
	adminMux.HandleFunc("/admin/film/update", f.UpdateFilm)
	adminMux.HandleFunc("/admin/film/delete", f.DeleteFilm)
	adminMux.HandleFunc("/admin/actor/update", a.UpdateActor)

	adminAuthHandler := auth.AdminAuthMiddleware(sm, adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminAuthHandler)
	siteMux.HandleFunc("/user/actor/add", a.AddActor)
	siteMux.HandleFunc("/user/film/add", f.AddFilm)
	siteMux.HandleFunc("/user/film/filmsList", f.GetAllFilms)
	siteMux.HandleFunc("/user/film/findFilms", f.FindFilms)
	siteMux.HandleFunc("/user/film/actorsListWithFilms", f.ActorsListWithFilms)

	siteMux.HandleFunc("/login", u.Login)
	siteMux.HandleFunc("/logout", u.Logout)
	siteMux.HandleFunc("/reg", u.Reg)

	http.Handle("/", auth.AuthMiddleware(sm, siteMux))

	http.ListenAndServe(":8080", nil)

}
