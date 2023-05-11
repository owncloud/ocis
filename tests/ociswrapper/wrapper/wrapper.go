package wrapper

import (
	"log"
	"net/http"
	"ociswrapper/common"
	"ociswrapper/ocis/config"
	"ociswrapper/wrapper/handlers"
)

func Start(port string) {
	defer common.Wg.Done()

	if port == "" {
		port = config.Get("port")
	}

	httpServer := &http.Server{
		Addr: ":" + port,
	}

	var mux = http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/config", handlers.SetEnvHandler)
	mux.HandleFunc("/rollback", handlers.RollbackHandler)

	httpServer.Handler = mux

	log.Printf("Starting server on port %s...", port)

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
