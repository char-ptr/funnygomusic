package databaser

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

type Indexer struct {
	db          *gorm.DB
	files       *sync.Map
	albums      []string
	albumsLock  sync.Mutex
	artists     []string
	artistsLock sync.Mutex
}

type IndexedArtist struct {
	gorm.Model
	Name string
	ID   string `gorm:"primaryKey"`
}
type IndexedAlbum struct {
	gorm.Model
	Name    string
	Genre   string
	Release string
	ID      string `gorm:"primaryKey"`
	Artist  string
}
type IndexedSong struct {
	gorm.Model
	Title  string
	Artist string
	Album  string
	Track  int
	File   string
	Length float64
	ID     string `gorm:"primaryKey"`
}

type DistinctIndexed struct {
	Artist string
	Album  string
	File   string
}

func NewIndexer(db *gorm.DB) *Indexer {
	//var distincts []DistinctIndexed
	rows, _ := db.Raw("select file, iso.album, artist\nfrom indexed_songs iso\n         inner join (select album, min(track) as minTrack\n                     from indexed_songs\n                     group by album) t on iso.album = t.album and iso.track = t.minTrack\n").Rows()
	defer rows.Close()
	var artists []string
	var albums []string
	files := &sync.Map{}
	log.Println("get distincts")
	for rows.Next() {
		var v DistinctIndexed
		//cols, _ := rows.Columns()
		db.ScanRows(rows, &v)
		//log.Println("lol", cols)
		//log.Println("got distinct", v.file, len(v.file))
		artists = append(artists, v.Artist)
		albums = append(albums, v.Album)
		files.Store(v.File, true)
	}
	log.Println("got distincts")
	return &Indexer{
		db:          db,
		files:       files,
		albums:      albums,
		artists:     artists,
		artistsLock: sync.Mutex{},
		albumsLock:  sync.Mutex{},
	}
}
func (i *Indexer) IndexDirectory(dir string, ctx context.Context) {
	sema := semaphore.NewWeighted(60)
	waiter := sync.WaitGroup{}

	toCommitSong := sync.Map{}
	toCommitAlbum := sync.Map{}
	toCommitArtist := sync.Map{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		_, ok := i.files.Load(path)
		log.Println("path:", path, "ok:", ok)
		if !ok {
			waiter.Add(1)
			err := sema.Acquire(ctx, 1)
			if err != nil {
				log.Println("failed to acquire semaphore:", err)
			}
			go i.IndexFile(path, ctx, &toCommitSong, &toCommitAlbum, &toCommitArtist, sema, &waiter)
		} else {
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		log.Println("failed to walk directory:", err)
		return
	}
	log.Println("waiting for indexing to finish")
	waiter.Wait()
	log.Println("committing to db")
	i.db.Transaction(func(tx *gorm.DB) error {
		toCommitSong.Range(func(key, value any) bool {
			val := value.(IndexedSong)
			log.Printf("committing song:%#v \n", val)
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&val).Error
			if err != nil {
				fmt.Println("failed to commit song:", err)
			}
			return true
		})
		toCommitAlbum.Range(func(key, value any) bool {
			val := value.(IndexedAlbum)
			tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&val)
			return true
		})
		toCommitArtist.Range(func(key, value any) bool {
			val := value.(IndexedArtist)
			tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&val)
			return true
		})
		return nil
	})
	log.Println("done indexing directory")

}

func (i *Indexer) IndexFile(file string, ctx context.Context, songm *sync.Map, albm *sync.Map, artm *sync.Map, sema *semaphore.Weighted, waiter *sync.WaitGroup) {
	defer waiter.Done()
	defer sema.Release(1)
	data, err := FetchDataForFile(file, ctx)
	if err != nil {
		log.Println("failed to index file:", err)
		return
	}
	i.artistsLock.Lock()

	if slices.Contains(i.artists, data.Format.Tags.ArtID) == false {
		i.artists = append(i.artists, data.Format.Tags.ArtID)
		newArtist := IndexedArtist{
			Name: data.Format.Tags.Artist,
			ID:   data.Format.Tags.ArtID,
		}
		artm.Store(data.Format.Tags.ArtID, newArtist)

	}
	i.artistsLock.Unlock()
	i.albumsLock.Lock()
	if slices.Contains(i.albums, data.Format.Tags.AlbID) == false {
		i.albums = append(i.albums, data.Format.Tags.AlbID)
		newAlbum := IndexedAlbum{
			Name:    data.Format.Tags.Album,
			Genre:   data.Format.Tags.Genre,
			Release: data.Format.Tags.Date,
			ID:      data.Format.Tags.AlbID,
			Artist:  data.Format.Tags.ArtID,
		}
		albm.Store(data.Format.Tags.AlbID, newAlbum)

	}
	i.albumsLock.Unlock()
	parseDuration, err := strconv.ParseFloat(data.Format.Duration, 64)
	if err != nil {
		log.Println("failed to parse duration:", err)
		return
	}
	parsedTrack, _ := strconv.Atoi(data.Format.Tags.Track)
	newSong := IndexedSong{
		Title:  data.Format.Tags.Title,
		Artist: data.Format.Tags.ArtID,
		Album:  data.Format.Tags.AlbID,
		Track:  parsedTrack,
		File:   file,
		Length: parseDuration,
		ID:     data.Format.Tags.TrkID,
	}
	songm.Store(data.Format.Tags.TrkID, newSong)
	i.files.Store(file, true)

}

type RawProbeOutput struct {
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Track  string `json:"track"`
			Title  string `json:"title"`
			Artist string `json:"artist"`
			Album  string `json:"album"`
			TrkID  string `json:"musicbrainz_trackid"`
			AlbID  string `json:"musicbrainz_albumid"`
			ArtID  string `json:"musicbrainz_artistid"`
			Date   string `json:"date"`
			Genre  string `json:"genre"`
		} `json:"tags"`
	} `json:"format"`
}

func FetchDataForFile(file string, ctx context.Context) (RawProbeOutput, error) {
	//println("fetching data for file:", file)
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration:format_tags=track,title,artist,album,musicbrainz_trackid,musicbrainz_albumid,musicbrainz_artistid,date,genre",
		"-of", "json=c=1",
		file)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		println("err running ffprobe:", err)
		return RawProbeOutput{}, err
	}
	var raw = RawProbeOutput{}
	err = json.Unmarshal(out, &raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "err parsing ffprobe:", err, string(out))
		//println("err parsing ffprobe:", err)
		return RawProbeOutput{}, err
	}
	return raw, nil
}
