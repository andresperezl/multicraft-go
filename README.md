# multicraft-go

This is a Multicraft API client implementation in Go.

## Usage

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
	
}
```

Full list of Multicraft API functions can be found [here](https://www.multicraft.org/site/docs/api#6)

### Host specific tweaks

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

## TODO

- Add actual functions for each of the functions supported by Multicraft API
    - Validate the parameters, at least the types
- Create CLI implementation of the client