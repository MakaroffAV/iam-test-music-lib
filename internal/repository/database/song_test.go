package database

import (
	"errors"
	"github.com/google/uuid"
	"music_lib/internal/config"
	"testing"
)

var (
	testSong  = NewSong(config.NewDB().MustConn())
	testVerse = NewVerse(config.NewDB().MustConn())

	testSongID       = uuid.NewString()
	testSongYouTube  = "https://www.youtube.com/watch?v=E0ozmU9cJDg"
	testSongReleased = "1986-03-03"
	testSongAlbum    = "Muster of Puppets"
	testSongTitle    = "Muster of Puppets"
	testSongArtist   = "Metallica"
	testSongVerses   = []string{
		"End of passion play, crumbling away",
		"I'm your source of self-destruction",
	}
)

func TestSong_Create(t *testing.T) {
	_, err := testSong.Create(
		testVerse, testSongVerses, testSongID, testSongYouTube,
		testSongReleased, testSongAlbum, testSongTitle, testSongArtist,
	)
	if !errors.Is(err, nil) {
		t.Fatalf(
			"test failed: got: %v; want: %v", err, nil,
		)
	}
}

func TestSong_Update(t *testing.T) {
	_, err := testSong.Update(
		testSongID, "", "", "", "", "",
	)
	if !errors.Is(err, nil) {
		t.Fatalf(
			"test failed: got: %v; want: %v", err, nil,
		)
	}
}

func TestSong_Search(t *testing.T) {
	_, err := testSong.Search(
		5, 0, "", "", "", "", "", "", "",
	)
	if !errors.Is(err, nil) {
		t.Fatalf(
			"test failed: got: %v; want: %v", err, nil,
		)
	}
}

func TestSong_Delete(t *testing.T) {
	err := testSong.Delete(testSongID)
	if !errors.Is(err, nil) {
		t.Fatalf(
			"test failed: got: %v; want: %v", err, nil,
		)
	}
}
