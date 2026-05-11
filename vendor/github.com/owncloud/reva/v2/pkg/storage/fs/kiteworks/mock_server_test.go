// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/mock_server_test.go
package kiteworks_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks"
)

func newMockServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/rest/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			ID:    "u1",
			Name:  "Alice",
			Email: "alice@example.com",
			Quota: QuotaInfo{Allowed: 10737418240, Used: 1073741824},
		})
	})

	mux.HandleFunc("/rest/folders/top", func(w http.ResponseWriter, r *http.Request) {
		isShared := false
		parentID := "0"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []FileInfo{
				{ID: "f1", Name: "MyFiles", Type: "d", IsShared: &isShared},
				{ID: "f2", Name: "SharedFolder", Type: "d", IsShared: func() *bool { b := true; return &b }(), ParentID: &parentID},
			},
		})
	})

	mux.HandleFunc("/rest/folders/f1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		isShared := false
		json.NewEncoder(w).Encode(FileInfo{ID: "f1", Name: "MyFiles", Type: "d", IsShared: &isShared})
	})

	mux.HandleFunc("/rest/folders/f1/children", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DirectoryInfo{
			Files:   []FileInfo{{ID: "file1", Name: "hello.txt", Type: "f", Size: 5}},
			Folders: []FileInfo{},
		})
	})

	mux.HandleFunc("/rest/files/file1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FileInfo{ID: "file1", Name: "hello.txt", Type: "f", Size: 5})
	})

	mux.HandleFunc("/rest/files/file1/content", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})

	return httptest.NewServer(mux)
}
