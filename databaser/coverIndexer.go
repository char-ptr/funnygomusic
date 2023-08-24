package databaser

import (
	"context"
	"github.com/google/uuid"
	"github.com/kirsle/configdir"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var configPath string

var (
	ConfigFolders = [1]string{
		"artwork",
	}
)

func MakeConfigPath() string {
	if configPath == "" {
		configPath = configdir.LocalConfig("funnygomusic")
		err := configdir.MakePath(configPath)
		if err != nil {
			slog.Error("Failed to create config path", "error", err)
			return ""
		}
	}
	for _, v := range ConfigFolders {
		err := os.MkdirAll(filepath.Join(configPath, v), 0755)
		if err != nil {
			slog.Error("Failed to create config path", "error", err)
			return ""
		}

	}
	return configPath
}

func IndexFileArtwork(path string, id uuid.UUID, wg *sync.WaitGroup, shore *semaphore.Weighted) (savedTo string) {
	defer wg.Done()
	defer shore.Release(1)
	writeToFile := filepath.Join(MakeConfigPath(), "artwork")
	err := os.MkdirAll(writeToFile, 0755)
	if err != nil {
		slog.Error("Failed to create config path", "error", err)
		return
	}
	savedTo = filepath.Join(writeToFile, id.String()) + ".webp"
	exec.Command("ffmpeg", "-i", path, "-an", "-vcodec", "copy", "-f", "image2", savedTo).Run()
	return
}

type FileAndId struct {
	File string
	Id   uuid.UUID
}

func UpdateIndexedArtworks(db *gorm.DB) {
	alreadyIndexed := []string{"00.0.000"}
	artworkDir := filepath.Join(MakeConfigPath(), "artwork")
	err := os.MkdirAll(artworkDir, 0755)
	if err != nil {
		slog.Error("Failed to create config path", "error", err)
		return
	}
	filepath.WalkDir(artworkDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fpBase := filepath.Base(path)
		alreadyIndexed = append(alreadyIndexed, strings.TrimSuffix(fpBase, filepath.Ext(fpBase)))
		return nil
	})
	var filesToDo []FileAndId
	log.Println(alreadyIndexed)
	db.Raw("select distinct on (ia.id) iso.file, ia.id from indexed_albums as ia left join indexed_songs as iso on iso.album = ia.id and ia.id not in ? where iso.file IS NOT NULL", alreadyIndexed).Scan(&filesToDo)
	log.Println(filesToDo)
	waitgpr := sync.WaitGroup{}
	semaph := semaphore.NewWeighted(20)
	for _, v := range filesToDo {
		waitgpr.Add(1)
		semaph.Acquire(context.TODO(), 1)
		go IndexFileArtwork(v.File, v.Id, &waitgpr, semaph)
	}
	waitgpr.Wait()
}
