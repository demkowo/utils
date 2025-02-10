# utils

This library provides:
- **config**: Configuration management (e.g. loading environment variables into `conf`)
- **resp**: Standardized responses (`Err` and `Ok` types with `.JSON()` methods)
- **helpers (utils)**: Common functions for JWT token creation, password hashing, UUID/time parsing, and test mocks
- **models**: Data structures (e.g. `Account`, `APIKey`) and validation logic

## Installation

```bash
go get github.com/demkowo/utils
```

## Usage

### Initialize config:
```go
import "github.com/demkowo/utils/config"

func main() {
    cfg := config.Values.Get()
    // use cfg.JWTSecret, etc.
}
```

### Responses:
```go
import "github.com/demkowo/utils/resp"

// Error response
errResp := resp.Error(500, "something went wrong", []interface{}{"cause details"})
code, body := errResp.JSON()

// OK response
okResp := resp.New("success", []interface{}{"some", "data"})
code, body := okResp.JSON()
```

### Helpers:
```go
import (
    "github.com/demkowo/utils"
    "github.com/demkowo/utils/models"
)

func createToken() (string, error) {
    account := &models.Account{ /* fields */ }
    token, errResp := utils.H.AddJWTToken(account)
    if errResp != nil {
        return "", fmt.Errorf("error creating token: %v", errResp)
    }
    return token, nil
}
```

### Testing with Mocks:
```go
import "github.com/demkowo/utils"

func TestHelpersMock(t *testing.T) {
    // Start mocks
    utils.StartMock()
    defer utils.StopMock()

    // Set expected error or password
    utils.Var.SetExpectedError(map[string]error{"HashPassword": nil})
    // ...
}
```

## Components

`config`

Manages environment-based configuration (JWTSecret), providing a global Values variable to retrieve settings.

`resp`

Encapsulates standard JSON responses with Err and Ok structures, each having a .JSON() method for consistent HTTP status, code, and message formatting.

`helpers (utils)`

Provides:
- JWT token creation (`AddJWTToken`)
- Binding & validation (`BindJSON`, `ParseTime`, `ParseUUID`)
- Cryptographic helpers (`HashPassword`, `GetRandomBytes`)
- Mock versions (`hMock`) for easy testing

`models`

Defines data structures (`Account`, `APIKey`, etc.) with built-in validation methods, focusing on common authentication/authorization fields and constraints.