package wrapper

import (
	"fmt"
	"net/http"
	"ociswrapper/common"
	"ociswrapper/log"
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
	mux.HandleFunc("/k8s/config", handlers.K8sSetEnvHandler)
	mux.HandleFunc("/k8s/rollback", handlers.K8sRollbackHandler)
	mux.HandleFunc("/rollback", handlers.RollbackHandler)
	mux.HandleFunc("/command", handlers.CommandHandler)
	mux.HandleFunc("/stop", handlers.StopOcisHandler)
	mux.HandleFunc("/start", handlers.StartOcisHandler)
	mux.HandleFunc("/services/{service}", handlers.OcisServiceHandler)
	mux.HandleFunc("/services/rollback", handlers.RollbackServicesHandler)

	httpServer.Handler = mux

	log.Println(fmt.Sprintf("Starting ociswrapper on port %s...", port))

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
