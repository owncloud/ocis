package tree

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

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

		if isLockFile(lwev.Path) || isTrash(lwev.Path) {
			continue
		}

		switch {
		case strings.Contains(lwev.Event, "IN_CREATE"):
			go func() { _ = w.tree.Scan(lwev.Path, ActionCreate, false, false) }()
		case strings.Contains(lwev.Event, "IN_CLOSE_WRITE"):
			bytesWritten, err := strconv.Atoi(lwev.BytesWritten)
			if err == nil && bytesWritten > 0 {
				go func() { _ = w.tree.Scan(lwev.Path, ActionUpdate, false, true) }()
			}
		case strings.Contains(lwev.Event, "IN_MOVED_TO"):
			go func() {
				_ = w.tree.Scan(lwev.Path, ActionMove, false, true)
				_ = w.tree.WarmupIDCache(lwev.Path, false)
			}()
		}
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
