package storage

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nmktad/bumflix/internal/config"
)

var client *s3.Client

func New() (*s3.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read environment variables\n%v\n", err)
	}

	if client == nil {
		sdkcfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("Couldn't load AWS default configuration\n%v\n", err)
		}

		client = s3.NewFromConfig(sdkcfg, func(o *s3.Options) {
			o.Region = cfg.AWSDefaultRegion
			o.BaseEndpoint = aws.String(cfg.AWSS3Endpoint)
		})
	}

	return client, nil
}

func GetObject(ctx context.Context, objectKey string, bucketName string, file *os.File) (error) {
	c, err := New()
	if err != nil {
		return fmt.Errorf("error encountered when creating a client: %w", err)
	}

	res, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err !=  nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)

	return nil
}

func UploadDir(ctx context.Context, dirPath string, bucketName string, prefix string) (error) {
	c, err := New()
	if err != nil {
		return fmt.Errorf("error encountered when creating a client: %w", err)
	}

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		rel, err := filepath.Rel(dirPath, path)

		if err != nil {
			return err
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}
		defer file.Close()

		key := filepath.ToSlash(filepath.Join(prefix, rel))

		_, err = c.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			return fmt.Errorf("uploading %s: %w", key, err)
		}

		return nil
	})
}
