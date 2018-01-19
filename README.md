# PDX, City of Donuts

Who doesn't love a fresh donut? No-one I know.

Here's [the easiest way to find 'em](https://rjz.github.io/pdxdonuts).

## Testimonials

> This makes me want a donut right now.
>
> -- @mikaelsnavy

## Install

With the go toolchain [installed and configured][install-golang], install
dependencies as usual...

```ShellSession
$ go get ./...
```

And just add API keys!

```ShellSession
# https://developers.google.com/places/web-service/get-api-key
$ export GOOGLE_API_KEY='<your api key here>'

# https://www.mapbox.com/help/how-access-tokens-work/
$ export MAPBOX_ACCESS_TOKEN='<your access token here>'

# Build your own donut map
$ go run main.go \
    -keyword donut \
    -type 'restaurant|bakery' \
    -location 'Portland, OR'
```

[install-golang]: https://golang.org/doc/install#testing
