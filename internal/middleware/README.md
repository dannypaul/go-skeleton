## Correlation Id middleware

Generally it is difficult to correlate the logs generated by a request.

The idea of the `X-Correlation-ID` header is that a client can create some random ID and pass it to the server. The server then include that ID in every log statement that it creates. If a client receives an error it can include the ID in a bug report, allowing the server operator to look up the corresponding log statements.

The correlation ID is also set by the server in the response header `X-Correlation-ID`

## Auth middleware

The `Authorization` header uses the following syntax:
`Authorization: <type> <credentials>`

When `type` is:
* `Basic`: credentials should contain the base64 encoded of `<client-id>:<client-secret>`
* `Bearer`: credentials should contain a JSON Web Token (JWT)

The auth middleware extracts the `Authorization` header from every request and returns HTTP status code `401` if it matches the following criteria:
* If it contains an invalid `type`
* If it contains an invalid `client-id` and/or `client-secret`. It does so why verifying if the `client-id` and `client-secret` pair is persisted in the `apikeys` collection
* If the JWT token is expired
* If the header, payload or signature of the JWT token is tampered

If the `Authorization` header contains a valid `<type>` and `<crendentials>`, the middleware adds the authenticated user information to the request `context`. 

The functions that handle the business logic are responsible to validate if the authenticated user has the necessary permissions to execute it. The `iam` package has a utility method to aid the business logic functions in authorization. 
```go 
func VerifySession(ctx context.Context, hasAnyRole []Role) (Claims, error) {
    ...
}
```  
🚧 NOTE: Handling `Basic` type is work in progress 
