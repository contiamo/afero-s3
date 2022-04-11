package main

import (
	"context"
	"flag"
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
	log.SetOutput(os.Stderr)

	flag.StringVar(&region, "region", region, "AWS region")
	flag.StringVar(&bucket, "bucket", bucket, "S3 bucket")
	flag.StringVar(&key, "key", key, "S3 key")
	flag.StringVar(&output, "output", output, "Output file, \"-\" for stdout")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	// Initialize the file system
	fs := s3fs.NewFsFromClient(bucket, s3Client)

	log.Printf("Copying s://%s/%s to %s", bucket, key, output)

	// And do your thing
	src, err := fs.Open(key)
	if err != nil {
		log.Fatalf("unable to open file, %v", err)
	}
	defer src.Close()

	var out = os.Stdout
	if output != "-" {
		file, err := os.Create(output)
		if err != nil {
			log.Printf("unable to create output file, %v", err)
			return
		}
		defer file.Close()
		out = file
	}

	n, err := io.Copy(out, src)
	if err != nil {
		log.Printf("unable to copy file, %v", err)
		return
	}
	log.Printf("copied %d bytes", n)
}
