module github.com/owncloud/ocis-store

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/blevesearch/bleve v1.0.9
	github.com/blevesearch/cld2 v0.0.0-20200327141045-8b5f551d37f5 // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/golang/protobuf v1.4.1
	github.com/ikawaha/kagome.ipadic v1.1.2 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.2.2-0.20200527082518-5641fa4a4c8c
	github.com/owncloud/ocis-settings v0.0.0-20200629120229-69693c5f8f43
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	github.com/tebeka/snowball v0.4.2 // indirect
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/ugorji/go v1.1.4 // indirect
	go.opencensus.io v0.22.3
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
