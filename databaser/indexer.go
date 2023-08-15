package databaser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
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
	albums      []uuid.UUID
	albumsLock  sync.Mutex
	artists     []uuid.UUID
	artistsLock sync.Mutex
}

type IndexedArtist struct {
	gorm.Model
	Name string
	ID   uuid.UUID `gorm:"primaryKey"`
}
type IndexedAlbum struct {
	gorm.Model
	Name    string
	Genre   string
	Release string
	ID      uuid.UUID `gorm:"primaryKey"`
	Artist  uuid.UUID
}
type IndexedSong struct {
	gorm.Model
	Title  string
	Artist uuid.UUID
	Album  uuid.UUID
	Track  int
	File   string
	Length float64
	ID     uuid.UUID `gorm:"primaryKey"`
}

type DistinctIndexed struct {
	artist uuid.UUID
	album  uuid.UUID
	file   string
}

func NewIndexer(db *gorm.DB) *Indexer {
	var disincts []DistinctIndexed
	db.Raw("select distinct artist,album,file from indexed_songs").Scan(&disincts)
	var artists []uuid.UUID
	var albums []uuid.UUID
	files := &sync.Map{}
	log.Println("get distincts")
	for _, v := range disincts {
		artists = append(artists, v.artist)
		albums = append(albums, v.album)
		files.Store(v.file, true)
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
		//log.Println("path:", path, "ok:", ok)
		if !ok {
			waiter.Add(1)
			err := sema.Acquire(ctx, 1)
			if err != nil {
				log.Println("failed to acquire semaphore:", err)
			}
			go i.IndexFile(path, ctx, &toCommitSong, &toCommitAlbum, &toCommitArtist, sema, &waiter)
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
			tx.Create(&val)
			return true
		})
		toCommitAlbum.Range(func(key, value any) bool {
			val := value.(IndexedAlbum)
			tx.Create(&val)
			return true
		})
		toCommitArtist.Range(func(key, value any) bool {
			val := value.(IndexedArtist)
			tx.Create(&val)
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
		artm.Store(data.Format.Tags.ArtID.String(), newArtist)

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
		albm.Store(data.Format.Tags.AlbID.String(), newAlbum)

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
	songm.Store(data.Format.Tags.TrkID.String(), newSong)
	i.files.Store(file, true)

}

type RawProbeOutput struct {
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Track  string    `json:"track"`
			Title  string    `json:"title"`
			Artist string    `json:"artist"`
			Album  string    `json:"album"`
			TrkID  uuid.UUID `json:"musicbrainz_trackid"`
			AlbID  uuid.UUID `json:"musicbrainz_albumid"`
			ArtID  uuid.UUID `json:"musicbrainz_artistid"`
			Date   string    `json:"date"`
			Genre  string    `json:"genre"`
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
