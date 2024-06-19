package tree

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"
)

type GpfsFileAuditLoggingWatcher struct {
	tree *Tree
}

type lwe struct {
	Event        string
	Path         string
	BytesWritten string
}

func NewGpfsFileAuditLoggingWatcher(tree *Tree, auditLogFile string) (*GpfsFileAuditLoggingWatcher, error) {
	w := &GpfsFileAuditLoggingWatcher{
		tree: tree,
	}

	_, err := os.Stat(auditLogFile)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *GpfsFileAuditLoggingWatcher) Watch(path string) {
start:
	file, err := os.Open(path)
	if err != nil {
		// try again later
		time.Sleep(5 * time.Second)
		goto start
	}
	defer file.Close()

	// Seek to the end of the file
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		time.Sleep(5 * time.Second)
		goto start
	}

	reader := bufio.NewReader(file)
	ev := &lwe{}
	for {
		line, err := reader.ReadString('\n')
		switch err {
		case nil:
			err := json.Unmarshal([]byte(line), ev)
			if err != nil {
				continue
			}
			switch ev.Event {
			case "CREATE":
				go func() { _ = w.tree.Scan(ev.Path, false) }()
			case "CLOSE":
				bytesWritten, err := strconv.Atoi(ev.BytesWritten)
				if err == nil && bytesWritten > 0 {
					go func() { _ = w.tree.Scan(ev.Path, true) }()
				}
			case "RENAME":
				go func() {
					_ = w.tree.Scan(ev.Path, true)
					_ = w.tree.WarmupIDCache(ev.Path, false)
				}()
			}
		case io.EOF:
			time.Sleep(1 * time.Second)
		default:
			time.Sleep(5 * time.Second)
			goto start
		}
	}
}
