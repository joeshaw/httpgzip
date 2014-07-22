# httpgzip #

`httpgzip` is a Go package that provides an `http.Handle` wrapper that
transparently compresses the response payload with gzip if the
`Accept-Encoding: gzip` request header is provided.  It sets the
`Vary: Accept-Encoding` and `Content-Encoding: gzip` response headers.

## API ##

```go
mux := http.NewServeMux()
mux.Handle("/", IndexHandler)
mux.Handle("/login", LoginHandler)

http.ListenAndServe(":8000", httpgzip.Handler(mux))
```

That's it.  But it's also
[on godoc](http://godoc.org/github.com/joeshaw/httpgzip).

Without gzip:

```
$ curl -v http://localhost:8080/

> GET / HTTP/1.1
> User-Agent: curl/7.30.0
> Host: localhost:8000
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 22 Jul 2014 19:16:20 GMT
< Content-Length: 446
<
Lorem ipsum dolor sit amet, consectetur adipisicing
elit, sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
```

With gzip:
```
$ curl --compress -v http://localhost:8080/

> GET / HTTP/1.1
> User-Agent: curl/7.30.0
> Host: localhost:8000
> Accept: */*
> Accept-Encoding: deflate, gzip
>
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Content-Type: text/plain; charset=utf-8
< Vary: Accept-Encoding
< Date: Tue, 22 Jul 2014 19:17:41 GMT
< Content-Length: 292
<
Lorem ipsum dolor sit amet, consectetur adipisicing
elit, sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
```
