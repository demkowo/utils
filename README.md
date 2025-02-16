# utils

This repository provides utility packages for Go applications:

1. **`resp`**  
   A standardized way to create JSON responses with HTTP status codes.  
   - `resp.Err` represents an error response object.  
   - `resp.Ok` represents a success response object.  
   - Both types feature a `.JSON()` method that returns a status code and a JSON map.

2. **`helper`**  
   A collection of utility functions (and a mockable interface) for common tasks such as:  
   - Binding JSON requests with `gin`  
   - Exchanging OAuth2 codes for tokens  
   - Generating random bytes  
   - Hashing passwords (bcrypt)  
   - Parsing time and UUID values  
   - Signing JWT tokens  

   The `helper` package also supports a mock mode for testing:
   - `StartMock()` / `StopMock()` switch between mock and real behavior.
   - `AddMock(...)` registers function-specific mocks.

3. **`httpclient`**  
   A simple HTTP client interface (and a mockable implementation) for making HTTP requests:
   - Real client uses Go’s `net/http`.
   - Mock client allows you to inject predictable responses or errors.

## Installation

```bash
go get github.com/demkowo/utils
```

## Usage

### 1. Package `resp`:

Defines two types (`Err` and `Ok`) with consistent `.JSON()` methods for returning HTTP status codes and JSON bodies.

```go
import (
    "net/http"
    "github.com/demkowo/utils/resp"
)

// Creating an error response:
errResp := resp.Error(http.StatusInternalServerError, "Something went wrong", []interface{}{"cause details"})
statusCode, body := errResp.JSON()
// statusCode -> 500
// body -> map[string]interface{}{ "error": "Something went wrong", ... }

// Creating a success (OK) response:
okResp := resp.New(http.StatusOK, "Operation succeeded", []interface{}{"additional", "data"})
statusCode, body := okResp.JSON()
// statusCode -> 200
// body -> map[string]interface{}{ "message": "Operation succeeded", ... }

```

### 2. Package `helper`:

Provides:
- Binding & validation (`BindJSON`)
- Cryptographic helpers (`HashPassword`, `GetRandomBytes`)
- Time and UUID parsing (`ParseTime`, `ParseUUID`)
- JWT token signing (`TokenSignedString`)
- Mock infrastructure (`StartMock`, `StopMock`, `AddMock`) for easy testing.

```go
import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/demkowo/utils/helper"
    "github.com/google/uuid"
)

// Creating a new Helper instance:
h := helper.NewHelper()

// 2.1 BindJSON:
func handleSomeRequest(c *gin.Context) {
    var requestBody struct {
        // ...
    }
    if !h.BindJSON(c, &requestBody) {
        // Error response is already sent to the client.
        return
    }
    // proceed...
}

// 2.2 Generate random bytes:
bytes, errResp := h.GetRandomBytes(16)
if errResp != nil {
    // Handle the error response...
}

// 2.3 Hash a password:
hash, errResp := h.HashPassword("superSecret123")
if errResp != nil {
    // Handle the error response...
}

// 2.4 Parse time:
var parsed time.Time
if !h.ParseTime(c, time.RFC3339, "2023-02-03T10:00:00Z", &parsed) {
    // Error response is already sent.
    return
}

// 2.5 Parse UUID:
var id uuid.UUID
if !h.ParseUUID(c, "id", "550e8400-e29b-41d4-a716-446655440000", &id) {
    // Error response is already sent.
    return
}

```

### 3. Package `httpclient`:

Provides:
- Basic HTTP request methods (`Get`, `Post`, `Put`, `Patch`, `Delete`, `Head`, `Options`) using the standard Go `net/http` package.
- A unified interface for making HTTP calls.
- Mock infrastructure (`StartMock`, `StopMock`, `AddMock`) for simulating responses during testing.

```go
import (
    "github.com/demkowo/utils/httpclient"
)

// Real or mock HTTP client:
client := httpclient.NewClient()

// Perform requests:
resp, err := client.Get("https://api.example.com/data", nil)
if err != nil {
    // handle error
}
defer resp.Body.Close()

// or:
resp, err = client.Post("https://api.example.com/create", []byte(`{"key":"value"}`), map[string]string{"Content-Type":"application/json"})
if err != nil {
    // handle error
}
defer resp.Body.Close()

```

### Testing with Mocks:
```go
import (
    "testing"
    "github.com/demkowo/utils/helper"
    "github.com/demkowo/utils/httpclient"
)

func TestSomething(t *testing.T) {
    // Mock the helper:
    helper.StartMock()
    defer helper.StopMock()

    // For example, to mock HashPassword:
    helper.AddMock(helper.Mock{
        Test: "TestSomething",
        Error: map[string]error{
            "HashPassword": nil, // or an actual error to simulate a failure
        },
        Password: "mockedHashedPassword",
    })

    // Mock the HTTP client:
    httpclient.StartMock()
    defer httpclient.StopMock()
    httpclient.AddMock(httpclient.Mock{
        Test: "TestSomething",
        Error: map[string]error{
            "Get":  nil, // or provide an error to simulate a GET failure
        },
        Response: http.Response{ /* ... */ },
    })

    // Now calls to httpclient.NewClient() and helper.NewHelper() will use mocks.
}
```
