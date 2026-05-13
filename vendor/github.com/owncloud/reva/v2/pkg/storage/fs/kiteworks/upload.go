// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/upload.go
package kiteworks

import (
	"io"
	"math"
)

const defaultChunkSize = 5 * 1024 * 1024 // 5 MB

// uploadFile performs a chunked upload of r into parentFolderID.
// chunkSize of 0 uses defaultChunkSize.
func uploadFile(c *Client, parentFolderID, filename string, size int64, r io.Reader, chunkSize int64) (*FileInfo, error) {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	numChunks := 1
	if size > 0 {
		numChunks = int(math.Ceil(float64(size) / float64(chunkSize)))
	}

	result, err := c.InitiateUpload(parentFolderID, filename, size, numChunks)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	var fi *FileInfo
	for i := 0; i < numChunks; i++ {
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return nil, err
		}
		isLast := i == numChunks-1
		fi, err = c.UploadChunk(result.URI, filename, buf[:n], i, isLast)
		if err != nil {
			return nil, err
		}
	}
	return fi, nil
}
