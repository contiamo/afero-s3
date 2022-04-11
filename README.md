# S3 Backend for Afero

![Build](https://github.com/contiamo/afero-s3/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/contiamo/afero-s3)](https://goreportcard.com/report/contiamo/afero-s3)
[![GoDoc](https://godoc.org/github.com/contiamo/afero-s3?status.svg)](https://godoc.org/github.com/contiamo/afero-s3)

## About

It provides an [afero filesystem](https://github.com/spf13/afero/) implementation of an [S3](https://aws.amazon.com/s3/) backend.

There are some other alternatives, but this implementation focuses on efficient memory usage by streaming the file download and uploads.

We are open to any improvement through issues or pull-request that might lead to a better implementation or even better testing.

## Known limitations

- File appending / seeking for write is not supported because S3 doesn't support it, it could be simulated by rewriting entire files.
- Chtimes is not supported because S3 doesn't support it, it could be simulated through metadata.
- Chmod support is very limited

## How to use

Note: Errors handling is skipped for brevity but a complete example is provided in the [`example` folder](./example/main.go)

```golang
package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	s3fs "github.com/contiamo/afero-s3"
)

var (
	region = "us-west-2"
	bucket = "my-bucket"
	key    = "/path/to/file.txt"
	output = "-"
)

func main() {
	ctx := context.Background()

  // initialize the S3 client
	cfg, _ := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	s3Client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	// Initialize the file system
	fs := s3fs.NewFsFromClient(bucket, s3Client)

	// And do your thing
	src, _ := fs.Open(key)
	defer src.Close()

	var out = os.Stdout
	n, _ := io.Copy(out, src)

	log.Printf("copied %d bytes", n)
```

## Development

The project uses [`Taskfile`](https://taskfile.dev/#/) to orchestrate the local development flow.

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

Install `task`, and then use

```sh
# see the available dev tasks
task --list
```

To run the test suite:

```sh
task test
```

To run the example code:

```sh
task run-example -- --help
```

## Thanks

The initial code (which was massively rewritten) comes from:

- [fclairamb's fork](https://github.com/fclairamb/afero-s3)
- Which comes from [wreulicke's fork](https://github.com/wreulicke/afero-s3)
- Itself forked from [aviau's fork](https://github.com/aviau/).
- Initially proposed as [an afero PR](https://github.com/spf13/afero/pull/90) by [rgarcia](https://github.com/rgarcia) and updated by [aviau](https://github.com/aviau).
