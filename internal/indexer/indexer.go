package indexer

import (
	"os"
	"path/filepath"
	"log/slog"

	"fs-chunker/internal/types"
)

// Crawl traverses a given directory and sends file paths via chan
func Crawl(dir string, jobs chan<-types.FileInfo) error {
	defer close(jobs)
	
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// skip any dirs
			if d.IsDir() {
				return nil
			}

			info, _ := d.Info()
			sizeBytes := info.Size() 

			slog.Info("found file", 
				slog.String("file", path),
				slog.Int64("sizeb", sizeBytes),
			)
			
			jobs <- types.FileInfo{
				Path: path,
				SizeB: sizeBytes,
			}
			
			return nil
	})
}