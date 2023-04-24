// Package metrics provides prometheus metrics for the data managers..
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DownloadsActive is the number of active downloads
	DownloadsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "reva_download_active",
		Help: "Number of active downloads",
	})
	// UploadsActive is the number of active uploads
	UploadsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "reva_upload_active",
		Help: "Number of active uploads",
	})
)
