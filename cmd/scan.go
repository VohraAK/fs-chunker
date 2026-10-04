/*
Copyright © 2026 Abdullah K Vohra <abdullah.k.vohra@gmail.com>
*/
package cmd

import (
	"os"
	"sync"
	"runtime"
	"fmt"
	"log/slog"
	"path/filepath"
	"encoding/json"

	"github.com/spf13/cobra"

	"fs-chunker/internal/types"
	"fs-chunker/internal/indexer"
	"fs-chunker/internal/chunker"
)

var dirToScan string
var workers int

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a directory and chunk its files",
	
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("scan start",
			slog.String("directory", dirToScan),
		)

		// standardise into absolute paths
		absDir, err := filepath.Abs(dirToScan)
		if err != nil {
			slog.Error("failed to convert path to absolute",
				slog.String("error", err.Error()),
			)
		}

		filepathChan := make(chan types.FileInfo, types.JobsChanSize)
		resultsChan := make(chan types.FileInfo, types.JobsChanSize)

		// launch a goroutine for the crawler
		go indexer.Crawl(absDir, filepathChan)
		
		// init worker pool
		wg := &sync.WaitGroup{}

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go chunker.Worker(filepathChan, resultsChan, wg)
		}

		// a seperate goroutine should wait for the workers, and close the channel ONLY when all are done
		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		// collect the results and process them
		fileManifest := types.Manifest{}

		for fileInfo := range resultsChan {
			fileManifest.Files = append(fileManifest.Files, fileInfo)
		}

		// encode the manifest to json
		jsonData, err := json.MarshalIndent(fileManifest, "", "  ")
		if err != nil {
			slog.Error("failed to generate JSON", 
				slog.String("error", err.Error()),
			)
		}

		// print the final manifest
		fmt.Println(string(jsonData))
		
	},
}


func init() {
	rootCmd.AddCommand(scanCmd)

	defaultWorkers := runtime.NumCPU()
	
	// define the dir flag
	scanCmd.Flags().StringVarP(&dirToScan, "dir", "d", "", "Target directory to scan (required)")
	err := scanCmd.MarkFlagRequired("dir")

	// define the workers flag (non-required)
	scanCmd.Flags().IntVarP(&workers, "workers", "w", defaultWorkers, "Number of workers")

	if err != nil {
		slog.Error("failed to mark dir flag as required",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// optional: an output flag for the manifest
}
