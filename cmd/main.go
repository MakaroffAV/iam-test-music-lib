package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"music_lib/internal/api"
	"music_lib/internal/config"
	"music_lib/internal/repository/database"
	"music_lib/internal/service"
	"net/http"
)

func main() {
	d := config.NewDB().MustConn()
	defer func() {
		_ = d.Close()
	}()

	r := chi.NewRouter()
	s := api.NewServer(
		service.NewSong(
			database.NewSong(d),
			database.NewVerse(d),
		),
		service.NewVerse(
			database.NewSong(d),
			database.NewVerse(d),
		),
	)

	h := api.HandlerFromMuxWithBaseURL(s, r, "/api")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Fatalln(
			err, "сервер неожиданно завершил работу",
		)
	}

}
