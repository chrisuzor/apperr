# apperr

Application errors with a category, and nothing else.

A service returns a code and a message a client may see. The edge that talks to the client maps the code to a status. Workers and other non-HTTP callers use the same error and never import `net/http`.

```bash
go get github.com/chrisuzor/apperr
```

Go 1.23 or newer.

## Errors

```go
question, err := store.GetQuestion(ctx, id)
if err != nil {
    return nil, apperr.Replace(err, sql.ErrNoRows,
        apperr.New(apperr.CodeNotFound, "question not found"))
}
```

`Replace` uses `errors.Is`. A nil error, a nil sentinel, or a nil replacement is left unchanged.

`Wrap` keeps the cause for logs:

```go
return apperr.Wrap(apperr.CodeInternal, err, "load question")
```

`CodeOf` and `Message` read the outermost `*apperr.Error` through `%w` wrappers. An ordinary error classifies as `internal`.

| Code | Meaning |
| --- | --- |
| `invalid` | The request is malformed |
| `not_found` | The resource does not exist |
| `forbidden` | The caller may not do this |
| `conflict` | The write clashes with current state |
| `internal` | Anything else, including an unclassified error |

## HTTP

`httpx` maps a code to a status and a client-safe message. It does not write a body.

| Code | Status |
| --- | --- |
| `invalid` | 400 |
| `not_found` | 404 |
| `forbidden` | 403 |
| `conflict` | 409 |
| `internal`, or any other error | 500 |

`PublicMessage` returns the stored message for the four client codes. Internal errors return `internal error`. The cause remains on the error value for logs.

```go
status := httpx.Status(err)
http.Error(w, httpx.PublicMessage(err), status)
```

## Gin

`apperrgin` is a separate module so the core does not depend on Gin.

```bash
go get github.com/chrisuzor/apperr/apperrgin
```

`Write` classifies the error and calls a renderer you supply. A nil error does not write.

```go
apperrgin.Write(c, err, func(c *gin.Context, status int, code, message string) {
    c.JSON(status, gin.H{
        "status":     http.StatusText(status),
        "message":    message,
        "error_code": code,
    })
})
```

## Tests

```bash
go test ./...
go test ./apperrgin/
```

`apperrgin` uses a `replace` directive so it builds against this repository before a version is tagged.
