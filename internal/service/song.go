package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"music_lib/internal/domain"
	repo "music_lib/internal/repository"
	"strings"
)

var (
	ErrSongNotFound = errors.New(
		"Песня не найдена",
	)
)

type Song struct {
	repoSong  repo.Song
	repoVerse repo.Verse
}

func NewSong(repoSong repo.Song, repoVerse repo.Verse) Song {
	return Song{
		repoSong:  repoSong,
		repoVerse: repoVerse,
	}
}

func (s Song) Delete(id string) error {
	// Удаляем песню из базы данных, по
	// внешнему ключу удалится и текст песни

	// Проверяем есть
	// ли песня в базе данных
	_, err := s.repoSong.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSongNotFound
		}
		return err
	}

	// Если песня
	// есть в базе - просто удаляем ее
	return s.repoSong.Delete(id)
}

func (s Song) Update(
	id, youtube, released, album, title, artist string) (domain.Song, error) {

	// Проверяем есть
	// ли песня в базе данных
	_, err := s.repoSong.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Song{}, ErrSongNotFound
		}
		return domain.Song{}, err
	}

	// Пишем в базу данных
	// обновленные данные песни
	return s.repoSong.Update(
		id, youtube, released, album, title, artist,
	)
}

func (s Song) Create(
	text, youtube, released, album, title, artist string) (domain.Song, error) {

	// Пишем в базу данных
	// информацию о новой песне
	//
	// В ТЗ приведен пример записи новой песни,
	// там в теле post передается текст песни как строка.
	//
	// Опять же по ТЗ есть условие, что мы должны разделять
	//куплеты песни, так что будем просить  от пользователя
	// последовательсноть '$$' как способ разделения куплетов
	//
	return s.repoSong.Create(
		s.repoVerse, strings.Split(text, "$$"),
		uuid.NewString(), youtube, released, album, title, artist,
	)
}

func (s Song) Search(
	limit, offset int,
	text, album, title, artist, youtube, releasedFrom, releasedTo string) ([]domain.Song, error) {

	// Вытаскиваем из
	// базы результаты поиска
	return s.repoSong.Search(
		limit, offset,
		text, album, title, artist, youtube, releasedFrom, releasedTo,
	)
}
