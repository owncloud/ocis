# meshdirectory-web

![image](demo/preview.png)

A web frontend SPA for the Science Mesh Directory service written in Vue 3.

## Development

### Project setup
```
yarn install
```

### Compiles and hot-reloads for development
```
yarn serve
```

### Compiles and minifies for production
```
yarn build
```

### Lints and fixes files
```
yarn lint
```

## Usage

To use this frontend in your GOlang projects, run the following:
```
go get github.com/sciencemesh/meshdirectory-web
```

And serve the SPA distribution in your HTTP handler, e.g.:
```go
package mypackage
import (
	"net/http"
	"log"
	"github.com/sciencemesh/meshdirectory-web"
)

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    return ServeMeshDirectorySPA(w, r)
	}
}

func main() {
    http.Handle("/", Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Finally try to access the Mesh Directory frontend by opening the following url in your browser:
```
http://localhost:8080/?token=xyz&providerDomain=cesnet.cz
```

## Credits

- Custom GeoJSON map data was generated using the [GeoJSON Maps](https://geojson-maps.ash.ms/) service.

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
