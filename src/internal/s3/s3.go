package S3

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

func New(bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
		),
	)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.Region = os.Getenv("AWS_DEFAULT_REGION")
			o.BaseEndpoint = aws.String(os.Getenv("AWS_S3_LOCALSTACK_ENDPOINT"))
		}),
		Bucket: bucket,
	}, nil
}

func (u *S3Client) UploadFile(localPath, remoteKey, contentType string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", localPath, cerr)
		}
	}()

	_, err = u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &u.Bucket,
		Key:         &remoteKey,
		Body:        file,
		ContentType: &contentType,
	})

	return err
}

func (u *S3Client) UploadDir(localPath, remotePrefix string) error {
	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
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
		key := filepath.ToSlash(filepath.Join(remotePrefix, rel)) // always use `/` for S3

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", path, err)
		}

		_, err = u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      &u.Bucket,
			Key:         aws.String(key),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			if cerr := file.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
			}

			return fmt.Errorf("upload failed for %s: %w", key, err)
		}

		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
		}

		fmt.Println("✅ Uploaded:", key)
		return nil
	})
}

func (u *S3Client) GetPresignedURL(objectName string) (string, error) {
	presignClient := s3.NewPresignClient(u.Client)
	presigner := Presigner{PresignClient: presignClient}

	params := url.Values{}

	if strings.HasSuffix(objectName, ".m3u8") {
		params.Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.HasSuffix(objectName, ".ts") {
		params.Set("Content-Type", "video/mp2t")
	} else {
		params.Set("Content-Type", "application/octet-stream")
	}

	presignedGetRequest, err := Presigner.GetObject(presigner, context.TODO(), os.Getenv("AWS_S3_BUCKET_NAME"), objectName, 60)
	if err != nil {
		panic(err)
	}

	return presignedGetRequest.URL, nil
}
