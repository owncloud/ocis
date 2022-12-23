package service

import (
	"fmt"
	"github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/r3labs/sse/v2"
	"net/http"
	"time"
)

type SSE struct{}

func ServeSSE(r chi.Router) {
	server := sse.New()
	stream := server.CreateStream("messages")
	stream.AutoReplay = false

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		u, ok := ctx.ContextGetUser(r.Context())
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go func() {
			for range time.Tick(time.Second * 4) {
				t := time.Now()
				server.Publish("messages", &sse.Event{
					Data: []byte(fmt.Sprintf("[%s] Hello %s, new push notification from server!", t.Format("2006-01-02 15:04:05"), u.Username)),
				})
			}
		}()

		server.ServeHTTP(w, r)
	})

}
