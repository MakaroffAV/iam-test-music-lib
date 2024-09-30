package database

import (
	"errors"
	"fmt"
	"testing"
)

func TestVerse_GetBySongID(t *testing.T) {
	r, err := testVerse.GetBySongID(testSongID, 5, 0)
	fmt.Println(r)
	if !errors.Is(err, nil) {
		t.Fatalf(
			"test failed: got: %v; want: %v", err, nil,
		)
	}
}
