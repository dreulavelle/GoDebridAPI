# GoDebridAPI

A Go wrapper for the Real-Debrid API.

## Installation

```bash
go get github.com/dreulavelle/GoDebridAPI
```

## Example

```go
import "github.com/dreulavelle/GoDebridAPI/api"

// Initialize client
client := api.HttpClient(api.GetApiKey())

// Get user details
user, err := client.RdGetUser()
```
