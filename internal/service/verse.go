package service

import (
	"database/sql"
	"errors"
	"music_lib/internal/domain"
	repo "music_lib/internal/repository"
)

type Verse struct {
	repoSong  repo.Song
	repoVerse repo.Verse
}

func NewVerse(repoSong repo.Song, repoVerse repo.Verse) Verse {
	return Verse{
		repoSong:  repoSong,
		repoVerse: repoVerse,
	}
}

func (v Verse) GetBySongID(limit, offset int, songID string) ([]domain.Verse, error) {

	// проверяем, что песня существует
	_, err := v.repoSong.GetByID(songID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSongNotFound
		}
		return nil, err
	}

	// Вытаскиваем из
	// базы данных результаты поиска
	return v.repoVerse.GetBySongID(songID, limit, offset)
}
