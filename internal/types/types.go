package types

// shared package which defines the shape of structures and the manifest

// a chunk is a sequence of bytes with a size and associated hash
// workers will send the chunk infos into channel
type Chunk struct {		
	Sequence int	`json:"sequence"`
	Hash string		`json:"hash"`
	Size int64		`json:"size"`
}

// a file will have a filename, filesize, and a collection of chunks
type FileInfo struct {
	Path string		`json:"filepath"`
	SizeB int64 	`json:"sizeb"`
	Chunks []Chunk 	`json:"chunks"`
}

// a manifest is a collection of FileInfos 
type Manifest struct {
	Files []FileInfo `json:"files"`
}

const ChunkSize = 4_194_304		// 4MB chunks by default
const JobsChanSize = 100