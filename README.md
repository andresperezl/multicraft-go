# multicraft-go
[![Go](https://github.com/andresperezl/multicraft-go/actions/workflows/go.yml/badge.svg)](https://github.com/andresperezl/multicraft-go/actions/workflows/go.yml)[![Coverage Status](https://coveralls.io/repos/github/andresperezl/multicraft-go/badge.svg)](https://coveralls.io/github/andresperezl/multicraft-go)[![Go Reference](https://pkg.go.dev/badge/github.com/andresperezl/multicraft-go.svg)](https://pkg.go.dev/github.com/andresperezl/multicraft-go)

This is a Multicraft API client implementation in Go.

# Usage

First get the dependency in your project with:

```shell
go get -u github.com/andresperezl/multicraft-go
```

The in your code, you could do something like:

```go
import (
	"github.com/andresperezl/multicraft-go/client"
)

func main() {
	mcClient := client.New("http://example.com/multicraft/api.go", "demo", "#6nh%tX=ot$sBX")

	resp, err := mcClient.Do("getServerStatus", map[string]string{
		"id": "1",
	})
	if err != nil {
		panic(err)
	}
	// Do something with the response here
}
```

_Full list of Multicraft API functions can be found [here](https://www.multicraft.org/site/docs/api#6)_

All the responses from the Multicraft API client are in the form:
```go
type MulticraftResponse struct {
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}
```

## Host specific tweaks

`MulticraftAPIClient` uses `*resty.Client` as it underlying HTTP client, which 
allows you to add any specific requirements that your hosting provider may have.
For example, if you are using PebbleHost you need to add some [specific headers](https://help.pebblehost.com/en/article/using-the-pebblehost-game-panel-api-mv0hk4/)):

```go
mc := New(
  "https://panel.pebblehost.com/api.php",
  "youremail@example.com",
  "yourApiKey")
mc.Client.SetHeader(
  "User-Agent",
  "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
mc.Client.SetHeader("Referer", "https://panel.pebblehost.com")
```

# TODO

- Add actual functions for each of the functions supported by Multicraft API
    - Validate the parameters, at least the types
    - Provide specific response types and not just `interface{}`
- Create CLI implementation of the client
