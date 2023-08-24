package databaser

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"log"
	"sort"
	"strings"
)

type SmolSongData struct {
	Title  string
	ID     string
	Artist string
}

func TryFindSong(str string, db *gorm.DB) *SmolSongData {
	var songs []SmolSongData
	newStr := "%" + strings.ReplaceAll(str, "%", "%%") + "%"

	db.Raw("select lower(indexed_songs.title) as title, indexed_songs.id, lower(indexed_artists.name) as artist from indexed_songs left join indexed_artists on indexed_artists.id = indexed_songs.artist where concat(lower(title),' ',lower(indexed_artists.name)) like lower(?)", newStr).Scan(&songs)

	log.Printf("Found %d songs", len(songs))

	var fuzzyFindSlice []string
	for _, song := range songs {
		fuzzyFindSlice = append(fuzzyFindSlice, song.Title+" "+song.Artist)
	}

	results := fuzzy.RankFind(strings.ToLower(str), fuzzyFindSlice)
	log.Printf("Found %d results: %#v\nold:%#v", len(results), results, fuzzyFindSlice)
	if len(results) == 0 {
		return nil
	}
	// return first result
	sort.Sort(results)
	return &songs[slices.Index(fuzzyFindSlice, results[0].Target)]
}
