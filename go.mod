module github.com/owncloud/ocis

go 1.17

require (
	github.com/CiscoM31/godata v1.0.4
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/logger/zerolog/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/kubernetes/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/mdns/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/v3 v3.6.1-0.20210924081004-8c39b1e1204d
	github.com/blevesearch/bleve/v2 v2.1.0
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/cs3org/go-cs3apis v0.0.0-20211007101428-6d142794ec11
	github.com/cs3org/reva v1.13.1-0.20211012070754-5dbb87c8e084
	github.com/disintegration/imaging v1.6.2
	github.com/glauth/glauth v1.1.3-0.20210729125545-b9aecdfcac31
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/render v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/iancoleman/strcase v0.2.0
	github.com/justinas/alice v1.2.0
	github.com/libregraph/lico v0.34.1-0.20210803054646-b584e0372224
	github.com/mennanov/fieldmask-utils v0.4.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484
	github.com/nmcclain/ldap v0.0.0-20210720162743-7f8d1e44eeba
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/owncloud/open-graph-api-go v0.0.0-20210511151655-57894f7d46fb
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.25.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/suture/v4 v4.0.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/yaegashi/msgraph.go v0.1.4
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.25.0
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/examples v0.0.0-20210802225658-edb9b3bc2266 // indirect
	google.golang.org/protobuf v1.27.1
	gotest.tools/v3 v3.0.3
	stash.kopano.io/kgol/rndm v1.1.1
)

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.8.1 // indirect
	github.com/RoaringBitmap/roaring v0.7.3 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.40.46 // indirect
	github.com/beevik/etree v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/blevesearch/bleve_index_api v1.0.1 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/mmap-go v1.0.2 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.0.1 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.1 // indirect
	github.com/blevesearch/vellum v1.0.5 // indirect
	github.com/blevesearch/zapx/v11 v11.2.2 // indirect
	github.com/blevesearch/zapx/v12 v12.2.2 // indirect
	github.com/blevesearch/zapx/v13 v13.2.2 // indirect
	github.com/blevesearch/zapx/v14 v14.2.2 // indirect
	github.com/blevesearch/zapx/v15 v15.2.2 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/bmizerany/pat v0.0.0-20170815010413-6226ea591a40 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/crewjam/httperr v0.2.0 // indirect
	github.com/crewjam/saml v0.4.3 // indirect
	github.com/cubewise-code/go-mime v0.0.0-20200519001935-8c5762b177d8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/eternnoir/gncp v0.0.0-20170707042257-c70df2d0cd68 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gdexlab/go-render v1.0.1 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-kit/log v0.1.0 // indirect
	github.com/go-ldap/ldap/v3 v3.4.1 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gomodule/redigo v1.8.5 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/schema v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-plugin v1.4.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.1 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/longsleep/go-metrics v0.0.0-20191013204616-cddea569b0ea // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.0.0-20201213122252-bcd7e1b9601e // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mendsley/gojwk v0.0.0-20141217222730-4d5ec6e58103 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mileusna/useragent v1.0.2 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/minio-go/v7 v7.0.14 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/nats-io/jwt v1.1.0 // indirect
	github.com/nats-io/nats.go v1.10.0 // indirect
	github.com/nats-io/nkeys v0.1.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pkg/xattr v0.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/pquerna/otp v1.3.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.28.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rickb777/date v1.12.4 // indirect
	github.com/rickb777/plural v1.2.0 // indirect
	github.com/rs/cors v1.8.0 // indirect
	github.com/rs/xid v1.3.0 // indirect
	github.com/russellhaering/goxmldsig v1.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sciencemesh/meshdirectory-web v1.0.4 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sony/gobreaker v0.4.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/steveyen/gtreap v0.1.0 // indirect
	github.com/studio-b12/gowebdav v0.0.0-20210917133250-a3a86976a1df // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tus/tusd v1.6.0 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.etcd.io/etcd/api/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210921065528-437939a70204 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	stash.kopano.io/kgol/kcc-go/v5 v5.0.1 // indirect
	stash.kopano.io/kgol/oidc-go v0.3.1 // indirect
)

// this is a transitive replace. See https://github.com/libregraph/lico/blob/master/go.mod#L38
replace github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
