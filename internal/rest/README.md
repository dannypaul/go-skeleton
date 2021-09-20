## Error response

All the functions that handle http requests use the following function in `internal/rest` package to respond to requests
```go
func EncodeRes(w http.ResponseWriter, r *http.Request, res interface{}, err error) {
...
}
```

When a request results in an error, the `EncodeRes` responds to the client with a body that follows the following syntax:
 ```json
{
	"errors": [{
		"code": "<error-code>",
		"message": "<error-message>"
	}],
	"correlationId": "<correlation-id>"
}
```

If the request had a valid string set in the `X-Correlation-ID` header, then it will be set in the `correlationId` key of the json response.

The `errors` key of the json response contains the list of error codes and their corresponding error messages of all the errors that occurred as part of this request.

All possible error codes, error messages, and their http status code can be found in the `internal/exception` package.

When an unknown error (one that is not defined in `internal/exception`) occurs `EncodeRes` responds with:
```json
{
	"errors": [{
		"code": "internalServerError",
		"message": "Internal Server Error"
	}],
	"correlationId": "<correlation-id>"
}
```

All the error responses are automatically logged by `func EncodeRes(w http.ResponseWriter, r *http.Request, res interface{}, err error)` at `info` level
