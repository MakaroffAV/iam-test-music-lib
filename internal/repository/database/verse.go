package database

import (
	"database/sql"
	"music_lib/internal/domain"
)

const (
	sqlVerseCreate = `
	insert into music_lib.verse(id, text, song_id, order_num) values ($1, $2, $3, $4)
	`

	sqlVerseGetBySongID = `
	select
	    t1.id,
	    t1.text,
	    t1.song_id,
	    t1.order_num
	from
		music_lib.verse as t1
	where
	    t1.song_id = $1
	limit
		$2
	offset
		$3
	`
)

type Verse struct {
	db *sql.DB
}

func NewVerse(db *sql.DB) Verse {
	return Verse{
		db: db,
	}
}

func (v Verse) Create(t *sql.Tx, songID, id, verse string, orderNum int) error {
	_, err := t.Exec(
		sqlVerseCreate,
		id, verse, songID, orderNum,
	)
	return err
}

func (v Verse) GetBySongID(songID string, limit, offset int) ([]domain.Verse, error) {
	var r []domain.Verse
	q, err := v.db.Query(
		sqlVerseGetBySongID,
		songID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = q.Close()
	}()

	for q.Next() {
		var c domain.Verse
		sErr := q.Scan(
			&c.ID, &c.Text, &c.SongID, &c.OrderNum,
		)
		if sErr != nil {
			return nil, sErr
		} else {
			r = append(r, c)
		}
	}
	return r, nil
}
