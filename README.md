# testcontainers-w-go

This is a "toy" example of how you could use [Testconatiners](https://testcontainers.com/) to move away from mocking
and use a real PostgreSQL dependency in your project's tests.

This repository was created as part of this blog post: https://www.anthony-turcios.dev/posts/testing-with-testcontainers-go/

## Requirements

* [Docker](https://docs.docker.com/engine/install/) 
* [go](https://go.dev/doc/install)

## Usage

```go
go test -v ./...
```