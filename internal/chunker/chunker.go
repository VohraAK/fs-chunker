package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"sync"

	"fs-chunker/internal/types"
)

// actual chunking logic
func processFile(file types.FileInfo) (types.FileInfo, error) {
	// read a file, 4MB at a time
	// maintain a sequence number for each chunk
	// hash each chunk
	// store seq, hash and size

	seqNo := 0 

	f, err := os.Open(file.Path)
	
	if err != nil {
		slog.Error("could not open file", 
			slog.String("Filepath", file.Path),
			slog.String("Error", err.Error()),
		)
		return file, err
	}
	
	defer f.Close()

	chunk := make([]byte, types.ChunkSize)
	chunks := []types.Chunk{}

	for {
		bytesRead, err := f.Read(chunk)

		if bytesRead > 0 {
			// hash the chunk 
			// IMP: hash the slice of chunk upto bytesRead!
			// not doing so might result in leftover data hashing
			chunkHash := sha256.Sum256(chunk[:bytesRead])
			hashString := hex.EncodeToString(chunkHash[:])
	
			chunks = append(chunks, types.Chunk{
				Sequence: seqNo,
				Hash: hashString,
				Size: int64(bytesRead),
			})
	
			// increment seqNo at end of loop
			seqNo += 1
		}
		
		if err != nil {
			// file ends?
			if err == io.EOF {
				break
			}
			
			slog.Error("could not read file", 
				slog.String("Filepath", file.Path),
				slog.String("Error", err.Error()),
			)
			return file, err
		}

	}
	
	file.Chunks = chunks
	return file, nil
}

func Worker(jobs <-chan types.FileInfo, results chan<- types.FileInfo, wg *sync.WaitGroup) {
	// worker pulls files from the channel
	// processes each file
	// sends the result to chan

	defer wg.Done()
	
	for file := range jobs {
		processedFile, err := processFile(file)
		
		if err != nil {
			slog.Error("failed to process file", 
				slog.String("error", err.Error()),
			)
			continue	// onto the next file
		}
		results <- processedFile
	}
}
