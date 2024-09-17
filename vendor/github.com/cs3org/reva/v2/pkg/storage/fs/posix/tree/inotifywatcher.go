package tree

import (
	"fmt"

	"github.com/pablodz/inotifywaitgo/inotifywaitgo"
)

type InotifyWatcher struct {
	tree *Tree
}

func NewInotifyWatcher(tree *Tree) *InotifyWatcher {
	return &InotifyWatcher{
		tree: tree,
	}
}

func (iw *InotifyWatcher) Watch(path string) {
	events := make(chan inotifywaitgo.FileEvent)
	errors := make(chan error)

	go inotifywaitgo.WatchPath(&inotifywaitgo.Settings{
		Dir:        path,
		FileEvents: events,
		ErrorChan:  errors,
		KillOthers: true,
		Options: &inotifywaitgo.Options{
			Recursive: true,
			Events: []inotifywaitgo.EVENT{
				inotifywaitgo.CREATE,
				inotifywaitgo.MOVED_TO,
				inotifywaitgo.CLOSE_WRITE,
				inotifywaitgo.DELETE,
			},
			Monitor: true,
		},
		Verbose: false,
	})

	for {
		select {
		case event := <-events:
			for _, e := range event.Events {
				if isLockFile(event.Filename) || isTrash(event.Filename) {
					continue
				}
				switch e {
				case inotifywaitgo.DELETE:
					go func() { _ = iw.tree.HandleFileDelete(event.Filename) }()
				case inotifywaitgo.CREATE:
					go func() { _ = iw.tree.Scan(event.Filename, ActionCreate, event.IsDir, false) }()
				case inotifywaitgo.MOVED_TO:
					go func() {
						_ = iw.tree.Scan(event.Filename, ActionMove, event.IsDir, true)
					}()
				case inotifywaitgo.CLOSE_WRITE:
					go func() { _ = iw.tree.Scan(event.Filename, ActionUpdate, event.IsDir, true) }()
				}
			}

		case err := <-errors:
			switch err.Error() {
			case inotifywaitgo.NOT_INSTALLED:
				panic("Error: inotifywait is not installed")
			case inotifywaitgo.INVALID_EVENT:
				// ignore
			default:
				fmt.Printf("Error: %s\n", err)
			}
		}
	}
}
