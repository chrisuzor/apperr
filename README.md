# apperr

Application errors with a category, and nothing else.

A service returns a code and a message a client may see. The edge that talks to the client maps the code to a status. Workers and other non-HTTP callers use the same error and never import `net/http`.

```bash
go get github.com/chrisuzor/apperr@v0.1.1
```

Go 1.21 or newer.

## Errors

```go
question, err := store.GetQuestion(ctx, id)
if err != nil {
    return nil, apperr.Replace(err, sql.ErrNoRows,
        apperr.New(apperr.CodeNotFound, "question not found"))
}
```

`Replace` uses `errors.Is`. A nil error, a nil sentinel, or a nil replacement is left unchanged. The sentinel is not kept on the chain. To keep it for logs, pass a replacement that already wraps it:

```go
apperr.Replace(err, sql.ErrNoRows, apperr.Wrap(apperr.CodeNotFound, err, "question not found"))
```

`Wrap` sets a new code and a new client message, and unwraps to the cause. That replaces any code already on the cause. `Annotate` adds context and leaves the code and client message alone:

```go
return apperr.Annotate(err, "load question")
```

`Wrap(nil)` and `Annotate(nil)` return nil, so `return apperr.Wrap(code, err, msg)` stays a success when `err` is nil. A failure path that wraps the wrong nil variable reports success the same way.

`CodeOf` and `Message` read the outermost `*apperr.Error` through `%w` wrappers. An ordinary error classifies as `internal`.

| Code | Meaning |
| --- | --- |
| `invalid` | The request is malformed |
| `unauthenticated` | The caller is not signed in |
| `forbidden` | The caller may not do this |
| `not_found` | The resource does not exist |
| `conflict` | The write clashes with current state |
| `rate_limited` | The caller should wait |
| `unavailable` | The service cannot handle the call right now |
| `internal` | Anything else, including an unclassified error |

## HTTP

`apperrhttp` maps a code to a status and a client-safe message. It does not write a body.

| Code | Status |
| --- | --- |
| `invalid` | 400 |
| `unauthenticated` | 401 |
| `forbidden` | 403 |
| `not_found` | 404 |
| `conflict` | 409 |
| `rate_limited` | 429 |
| `unavailable` | 503 |
| `internal`, or any other error | 500 |

`PublicMessage` returns the stored message for the codes in the table above, except `internal`. A 500 returns `internal error`, including when the code is not in the table. `Code` reports `internal` for those errors. `CodeOf` still returns the original code for logs. The cause remains on the error value.

```go
status := apperrhttp.Status(err)
http.Error(w, apperrhttp.PublicMessage(err), status)
```

## Gin

`apperrgin` is a separate module so the core does not depend on Gin. It needs Go 1.24.

```bash
go get github.com/chrisuzor/apperr/apperrgin@v0.1.1
```

`New` stores a renderer and panics if it is nil. `Write` classifies the error, records it with `c.Error`, stops the chain, and calls the renderer. A nil error does not write. `message` is safe for the client. `err` is the original error: use it for a typed cause, such as a validation error wrapped with `Wrap`. Handle a validation error that was never wrapped before calling `Write`, or `Write` classifies it as internal.

```go
write := apperrgin.New(func(c *gin.Context, status int, code, message string, err error) {
    c.JSON(status, gin.H{
        "status":     http.StatusText(status),
        "message":    message,
        "error_code": code,
    })
})

write.Write(c, err)
```

## Tests

This repository is a workspace of two modules.

```bash
go test -race ./...
go vet ./...
cd apperrgin && go test -race ./... && go vet ./...
```
