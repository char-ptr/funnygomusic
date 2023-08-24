package bot

import (
	"funnygomusic/databaser"
	"log"
)

func (b *Botter) MeiliUpdate() {

	sidx := b.Meili.Index("songs")
	var all_songs []databaser.IndexedSong
	b.Db.Raw("select iso.*, iaa.name as artist, ia.name as album from indexed_songs iso join public.indexed_albums ia on iso.album = ia.id join indexed_artists iaa on iso.artist = iaa.id where iso.id != ''").Scan(&all_songs)

	//for _, song := range all_songs {
	//	//log.Println(song.Title, song.Artist, song.Album)
	//}
	log.Println("songs", len(all_songs))

	tinfo, err := sidx.AddDocuments(all_songs)
	if err != nil {
		log.Println("failed to add songs to meili", err)
	}
	log.Println("meili request ok: ", tinfo.Status, tinfo.Type)
}
