# GoDebridAPI

A Go wrapper for the Real-Debrid API.

## Installation

```bash
go get github.com/dreulavelle/GoDebridAPI
```

## Usage

### Loading API Key

This can be done in one of the following ways. There is also an example to help get you started!

#### 1. Environment Variable

Set an environment variable in your shell:

```bash
export RD_API_KEY=your_real_debrid_api_key_here
```

In your Go code:

```go
apiKey := os.Getenv("RD_API_KEY")
client := GoDebridAPI.HttpClient(apiKey)
```

#### 2. Configuration File

Create a `config.json` file:

```json
{
  "api_key": "your_real_debrid_api_key_here"
}
```

In your Go code:

```go
var config map[string]string
file, _ := os.Open("config.json")
decoder := json.NewDecoder(file)
err := decoder.Decode(&config)
if err != nil {
  log.Fatalf("Error reading config: %v", err)
}
apiKey := config["api_key"]
client := GoDebridAPI.HttpClient(apiKey)
```

#### 3. Command-Line Argument

Run your Go program with the API key as an argument:

```bash
go run your_program.go --api-key your_real_debrid_api_key_here
```

In your Go code:

```go
var apiKey string
flag.StringVar(&apiKey, "api-key", "", "Your API Key")
flag.Parse()
client := GoDebridAPI.HttpClient(apiKey)
```

### Full Example

Here's a full example that uses the environment variable method to load the API key:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dreulavelle/GoDebridAPI" // Be sure to import this!
)

func main() {
	// Initialize the client
	apiKey := os.Getenv("RD_API_KEY")
	client := GoDebridAPI.HttpClient(apiKey)

	// Fetch user details
	user, err := client.RdGetUser()
	if err != nil {
		log.Fatalf("Error fetching user details: %v", err)
	}
	fmt.Printf("User Details: %+v\n", user)

	// Add more examples for other API calls
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](./LICENSE)
