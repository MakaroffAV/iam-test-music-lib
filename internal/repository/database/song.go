package database

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"music_lib/internal/domain"
	repo "music_lib/internal/repository"
)

const (
	sqlSongSearch = `
	select
	    t1.id,
	    t1.album,
	    t1.title,
	    t1.artist,
	    t1.youtube,
	    t1.released
	from
	    music_lib.song  as t1
	left join
		music_lib.verse as t2
	on
	    t1.id = t2.song_id
	where
	    t1.album like '%' || $1 || '%'
	  	and
	    t1.title like '%' || $2 || '%'
	    and
	    t1.artist like '%' || $3 || '%'
	    and
	    t1.youtube = COALESCE(NULLIF($4, ''), t1.youtube)
	    and
	    concat(t2.text, '$$') like '%' || $5 || '%'
		and
	    t1.released between COALESCE(NULLIF($6, '')::date, '1900-01-01'::date) and COALESCE(NULLIF($7, '')::date, '2100-01-01'::date)
	limit
		$8
	offset
		$9
	`

	sqlSongCreate = `
	insert into music_lib.song(id, youtube, released, album, title, artist) values ($1, $2, $3, $4, $5, $6)
	`

	sqlSongUpdate = `
	update music_lib.song
		set youtube = COALESCE(NULLIF($1, ''), youtube),
			released = COALESCE(NULLIF($2, '')::date, released),
			album = COALESCE(NULLIF($3, ''), album),
			title = COALESCE(NULLIF($4, ''), title),
			artist = COALESCE(NULLIF($5, ''), artist)
	where
	    id = $6
	`

	sqlSongDelete = `
	delete from music_lib.song where id = $1`

	sqlSongGetByID = `
	select id, youtube, released, album, title, artist from music_lib.song where id = $1
	`
)

type Song struct {
	db *sql.DB
}

func NewSong(db *sql.DB) Song {
	return Song{
		db: db,
	}
}

func (s Song) Delete(id string) error {
	_, err := s.db.Exec(
		sqlSongDelete, id,
	)
	return err
}

func (s Song) GetByID(id string) (domain.Song, error) {
	var r domain.Song
	err := s.db.QueryRow(
		sqlSongGetByID, id,
	).Scan(
		&r.ID, &r.Youtube, &r.Released, &r.Album, &r.Title, &r.Artist,
	)
	return r, err
}

func (s Song) Update(id, youtube, released, album, title, artist string) (domain.Song, error) {
	_, err := s.db.Exec(
		sqlSongUpdate,
		youtube, released,
		album, title, artist, id,
	)
	if err == nil {
		return s.GetByID(id)
	} else {
		return domain.Song{}, err
	}
}

func (s Song) Create(repoVerse repo.Verse, text []string, id, youtube, released, album, title, artist string) (domain.Song, error) {
	t, err := s.db.Begin()
	if err != nil {
		return domain.Song{}, err
	}

	// записываем
	// информацию о песне в БД
	if _, err := t.Exec(sqlSongCreate, id, youtube, released, album, title, artist); err != nil {
		if rErr := t.Rollback(); rErr != nil {
			log.Println(
				rErr,
				"произошла ошибка при отмене транзации",
			)
		}
		return domain.Song{}, err
	}

	// Записываем
	// куплеты песни в БД
	for i, v := range text {
		if cErr := repoVerse.Create(t, id, uuid.NewString(), v, i+1); cErr != nil {
			if rErr := t.Rollback(); rErr != nil {
				log.Println(
					rErr,
					"произошла ошибка при отмене транзации",
				)
			}
			return domain.Song{}, cErr
		}
	}

	// применяем изменения транзации и
	// отдаем пользователю созданную песню в БД
	if err := t.Commit(); err != nil {
		return domain.Song{}, err
	}
	return s.GetByID(id)
}

func (s Song) Search(limit, offset int, text, album, title, artist, youtube, releasedFrom, releasedTo string) ([]domain.Song, error) {
	var r []domain.Song

	q, err := s.db.Query(
		sqlSongSearch,
		album, title, artist, youtube, text,
		releasedFrom, releasedTo, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = q.Close()
	}()

	for q.Next() {
		var c domain.Song
		sErr := q.Scan(
			&c.ID, &c.Album, &c.Title,
			&c.Artist, &c.Youtube, &c.Released,
		)
		if sErr != nil {
			return nil, sErr
		} else {
			r = append(r, c)
		}
	}
	return r, nil
}
