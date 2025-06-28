package s3client

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = os.Getenv("AWS_DEFAULT_REGION")
		o.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL_S3"))
	}), nil
}

func UploadDir(localPath, remotePrefix, bucket string) error {
	client, err := New()
	if err != nil {
		return fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		var contentType string
		switch filepath.Ext(info.Name()) {
		case ".m3u8":
			contentType = "application/vnd.apple.mpegurl"
		case ".ts":
			contentType = "video/mp2t"
		default:
			contentType = "application/octet-stream"
		}

		rel, _ := filepath.Rel(localPath, path)
		key := filepath.ToSlash(filepath.Join(remotePrefix, rel))

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", path, err)
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
			}
		}()

		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			return fmt.Errorf("upload failed for %s: %w", key, err)
		}

		fmt.Println("✅ Uploaded:", key)
		return nil
	})
}

// func UploadDir(localPath, remotePrefix, bucket string) error {
// 	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil || info.IsDir() {
// 			return err
// 		}
//
// 		var contentType string
// 		switch filepath.Ext(info.Name()) {
// 		case ".m3u8":
// 			contentType = "application/vnd.apple.mpegurl"
// 		case ".ts":
// 			contentType = "video/mp2t"
// 		default:
// 			contentType = "application/octet-stream"
// 		}
//
// 		rel, _ := filepath.Rel(localPath, path)
// 		key := filepath.ToSlash(filepath.Join(remotePrefix, rel)) // always use `/` for S3
//
// 		file, err := os.Open(path)
// 		if err != nil {
// 			return fmt.Errorf("failed to open %s: %w", path, err)
// 		}
//
// 		client, err := New()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
// 			Bucket:      aws.String(bucket),
// 			Key:         aws.String(key),
// 			Body:        file,
// 			ContentType: aws.String(contentType),
// 		})
// 		if err != nil {
// 			if cerr := file.Close(); cerr != nil {
// 				fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
// 			}
//
// 			return fmt.Errorf("upload failed for %s: %w", key, err)
// 		}
//
// 		if cerr := file.Close(); cerr != nil {
// 			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
// 		}
//
// 		fmt.Println("✅ Uploaded:", key)
// 		return nil
// 	})
// }

func UploadFile(localPath, remoteKey, contentType, bucket string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", localPath, cerr)
		}
	}()

	client, err := New()
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         &remoteKey,
		Body:        file,
		ContentType: &contentType,
	})

	return err
}
