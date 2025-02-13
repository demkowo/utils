# utils

This library provides two main packages:
- **resp**: Standardized JSON responses (`Err` and `Ok` types with `.JSON()` methods)
- **helper**: Common utility functions and their mocks (for tasks such as parsing JSON, generating random bytes, hashing passwords, parsing UUID/time, etc.)

## Installation

```bash
go get github.com/demkowo/utils
```

## Usage

### Standardized Responses:
```go
import "github.com/demkowo/utils/resp"

// Error response
errResp := resp.Error(http.StatusInternalServerError, "Something went wrong", []interface{}{"cause details"})
statusCode, body := errResp.JSON() 
// Use statusCode as the HTTP status, body as the response JSON

// OK response
okResp := resp.New("Success message", []interface{}{"some", "data"})
statusCode, body := okResp.JSON()
```

### Helpers:
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/demkowo/utils/helper"
)

// Create a new helper instance
h := helper.NewHelper()

// 2.1 Bind JSON
func handleRequest(c *gin.Context) {
    var requestBody struct {
        // fields...
    }
    if !h.BindJSON(c, &requestBody) {
        return // JSON binding failed and the error response was already sent
    }
    // proceed with requestBody
}

// 2.2 Generate random bytes
randomBytes, errResp := h.GetRandomBytes(16)
if errResp != nil {
    // handle error (errResp contains an error response with code, causes, etc.)
}

// 2.3 Hash a password
hashed, errResp := h.HashPassword("somePassword123")
if errResp != nil {
    // handle error
}

// 2.4 Parse time
var parsedTime time.Time
if !h.ParseTime(c, time.RFC3339, "2021-09-11T09:00:00Z", &parsedTime) {
    return // error response was already sent
}

// 2.5 Parse UUID
var id uuid.UUID
if !h.ParseUUID(c, "myID", "550e8400-e29b-41d4-a716-446655440000", &id) {
    return // error response was already sent
}

// 2.6 Sign JWT tokens
tokenString, signErr := h.TokenSignedString(jwtToken, []byte("someJWTSecret"))
if signErr != nil {
    // handle signing error
}
```

### Testing with Mocks:
```go
import (
    "testing"
    "github.com/demkowo/utils/helper"
)

func TestSomething(t *testing.T) {
    // Enable mocks
    helper.StartMock()
    defer helper.StopMock()

    // Add a specific mock if needed
    helper.AddMock(helper.Mock{
        Test: "TestSomething",
        Error: map[string]error{
            "HashPassword": nil, // or an error to simulate failure
        },
        Password: "mockedHashedPassword",
    })

    // Now calls to helper.NewHelper().HashPassword will use the mock
}
```

## Components

`resp`

Defines two types (`Err` and `Ok`) with consistent `.JSON()` methods for returning HTTP status codes and JSON bodies.

`helper`

Provides:
- Binding & validation (`BindJSON`)
- Cryptographic helpers (`HashPassword`, `GetRandomBytes`)
- Time and UUID parsing (`ParseTime`, `ParseUUID`)
- JWT token signing (`TokenSignedString`)
- Mock infrastructure (`StartMock`, `StopMock`, `AddMock`) for easy testing.