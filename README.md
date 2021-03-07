# secretmanager

Golang library for managing configuration data from [Google Cloud's Secret Manager](https://cloud.google.com/solutions/secrets-management).

## Usage

```go
import "github.com/justindfuller/secretmanager"

var config struct {
  MySecret string `secretmanager:"MySecret"`
  AnotherSecret string `secretmanager:"AnotherSecret",version:"3"`
}

func main() {
  var c config
  if err := secretmanager.Parse(&c); err != nil {
    log.Fatal(err)
  }
  log.Print("Retrieved config from Google Secret Manager", c)
}
```

## Defaults

- `version` defaults to `latest`.
