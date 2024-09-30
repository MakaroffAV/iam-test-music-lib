package api

import (
	"encoding/json"
	"errors"
	"log"
	"music_lib/internal/service"
	"net/http"
)

type Server struct {
	serviceSong  service.Song
	serviceVerse service.Verse
}

func NewServer(serviceSong service.Song, serviceVerse service.Verse) Server {
	return Server{
		serviceSong:  serviceSong,
		serviceVerse: serviceVerse,
	}
}

func (_ Server) write(w http.ResponseWriter, b []byte, c int) error {
	w.WriteHeader(c)
	w.Header().Add(
		"Content-Type", "application/json",
	)
	if _, err := w.Write(b); err != nil {
		return err
	} else {
		return nil
	}
}

func (s Server) response(w http.ResponseWriter, b any, err error) error {
	var c int
	switch {
	default:
		{
			c = http.StatusInternalServerError
		}
	case errors.Is(err, nil):
		{
			c = http.StatusOK
		}
	case errors.Is(err, service.ErrSongNotFound):
		{
			c = http.StatusNotFound
		}
	}

	if c == http.StatusOK {
		m, mErr := json.Marshal(b)
		if mErr == nil {
			return s.write(w, m, c)
		} else {
			return s.write(w, nil, http.StatusInternalServerError)
		}
	} else {
		m, mErr := json.Marshal(
			struct {
				Reason string
			}{
				Reason: err.Error(),
			},
		)
		if mErr == nil {
			return s.write(w, m, c)
		} else {
			return s.write(w, nil, http.StatusInternalServerError)
		}
	}
}

// (GET /api/search)
func (s Server) SongSearch(w http.ResponseWriter, r *http.Request, params SongSearchParams) {
	var (
		l, o = 5, 0

		t, a, tl, ar, y, rf, rt string
	)

	// Заменяем стандартные значения
	// значениями из запроса пользователя
	if params.Limit != nil {
		l = *params.Limit
	}
	if params.Offset != nil {
		o = *params.Offset
	}

	//

	if params.Text != nil {
		t = *params.Text
	}
	if params.Album != nil {
		a = *params.Album
	}
	if params.Title != nil {
		tl = *params.Title
	}
	if params.Artist != nil {
		ar = *params.Artist
	}
	if params.Youtube != nil {
		y = *params.Youtube
	}
	if params.ReleasedTo != nil {
		rt = (*params.ReleasedTo).String()
	}
	if params.ReleasedFrom != nil {
		rf = (*params.ReleasedFrom).String()
	}

	// слой бизнес логики приложения
	v, err := s.serviceSong.Search(
		l, o, t, a, tl, ar, y, rf, rt,
	)
	if rErr := s.response(w, v, err); rErr != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}

// (POST /api/song/create)
func (s Server) SongCreate(w http.ResponseWriter, r *http.Request) {
	var (
		b SongCreateJSONRequestBody
		d = json.NewDecoder(r.Body)
	)

	// читаем тело запроса
	if err := d.Decode(&b); err != nil {
		if rErr := s.response(w, nil, err); rErr != nil {
			log.Println(
				err, "ошибка при отправке ответа клиенту",
			)
		}
	}

	// слой бизнес логики приложения
	v, err := s.serviceSong.Create(
		b.Text, b.Youtube,
		b.Released.String(), b.Album, b.Title, b.Artist,
	)
	if rErr := s.response(w, v, err); rErr != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}

// (DELETE /api/song/{songID}/delete)
func (s Server) SongDelete(w http.ResponseWriter, r *http.Request, songID SongID) {
	// слой бизнес логики приложения
	err := s.serviceSong.Delete(songID)
	if rErr := s.response(w, "ok", err); rErr != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}

// (PATCH /api/song/{songID}/update)
func (s Server) SongUpdate(w http.ResponseWriter, r *http.Request, songID SongID) {
	var (
		b SongUpdateJSONRequestBody
		d = json.NewDecoder(r.Body)

		// переменные для хранения
		// обновленных параметров песни
		y, rl, a, t, ar string
	)

	// читаем тело запроса
	if err := d.Decode(&b); err != nil {
		if rErr := s.response(w, nil, err); rErr != nil {
			log.Println(
				err, "ошибка при отправке ответа клиенту",
			)
		}
	}

	// дальше смотрим на переданные пераметры,
	// для не пустых - получаем значения в переменные
	if b.Album != nil {
		a = *b.Album
	}
	if b.Title != nil {
		t = *b.Title
	}
	if b.Artist != nil {
		ar = *b.Artist
	}
	if b.Youtube != nil {
		y = *b.Youtube
	}
	if b.Released != nil {
		rl = (*b.Released).String()
	}

	// слой бизнес логики приложения
	v, err := s.serviceSong.Update(songID, y, rl, a, t, ar)
	if rErr := s.response(w, v, err); rErr != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}

// (GET /api/song/{songID}/verse)
func (s Server) SongVerses(w http.ResponseWriter, r *http.Request, songID SongID, params SongVersesParams) {
	var (
		l, o = 5, 0
	)
	// Если параметры заданы
	// пользователем, используем их
	if params.Limit != nil {
		l = *params.Limit
	}
	if params.Offset != nil {
		o = *params.Offset
	}

	// Бизнес-логика приложения
	v, err := s.serviceVerse.GetBySongID(l, o, songID)
	if rErr := s.response(w, v, err); rErr != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}

// (GET /ping)
func (s Server) Ping(w http.ResponseWriter, r *http.Request) {
	if err := s.response(w, "ok", nil); err != nil {
		log.Println(
			err, "ошибка при отправке ответа клиенту",
		)
	}
}
