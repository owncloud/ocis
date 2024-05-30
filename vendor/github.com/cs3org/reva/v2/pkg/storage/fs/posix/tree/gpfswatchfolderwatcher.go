package tree

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/lookup"
	kafka "github.com/segmentio/kafka-go"
)

type GpfsWatchFolderWatcher struct {
	tree    *Tree
	brokers []string
}

func NewGpfsWatchFolderWatcher(tree *Tree, kafkaBrokers []string) (*GpfsWatchFolderWatcher, error) {
	return &GpfsWatchFolderWatcher{
		tree:    tree,
		brokers: kafkaBrokers,
	}, nil
}

func (w *GpfsWatchFolderWatcher) Watch(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: w.brokers,
		GroupID: "ocis-posixfs",
		Topic:   topic,
	})

	lwev := &lwe{}
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		err = json.Unmarshal(m.Value, lwev)
		if err != nil {
			continue
		}

		if strings.HasSuffix(lwev.Path, ".flock") || strings.HasSuffix(lwev.Path, ".mlock") {
			continue
		}

		switch {
		case strings.Contains(lwev.Event, "IN_CREATE"):
			go func() { _ = w.tree.Scan(lwev.Path, false) }()
		case strings.Contains(lwev.Event, "IN_CLOSE_WRITE"):
			bytesWritten, err := strconv.Atoi(lwev.BytesWritten)
			if err == nil && bytesWritten > 0 {
				go func() { _ = w.tree.Scan(lwev.Path, true) }()
			}
		case strings.Contains(lwev.Event, "IN_MOVED_TO"):
			go func() {
				_ = w.tree.Scan(lwev.Path, true)
				_ = w.tree.lookup.(*lookup.Lookup).WarmupIDCache(lwev.Path)
			}()
		}
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
