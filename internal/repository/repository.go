package repository

import (
	"database/sql"
	"music_lib/internal/domain"
)

type Song interface {
	Delete(id string) error
	GetByID(id string) (domain.Song, error)
	Update(id, youtube, released, album, title, artist string) (domain.Song, error)
	Create(repoVerse Verse, text []string, id, youtube, released, album, title, artist string) (domain.Song, error)
	Search(
		limit int, offset int,
		text, album, title, artist, youtube, releasedFrom, releasedTo string) ([]domain.Song, error)
}

type Verse interface {
	Create(t *sql.Tx, id, songID, verse string, orderNum int) error
	GetBySongID(songID string, limit, offset int) ([]domain.Verse, error)
}
