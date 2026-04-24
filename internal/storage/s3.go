package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/isa0-gh/reader/internal/config"
)

type S3Client struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	cdnPrefix string
}

func NewS3Client(cfg *config.Config) *S3Client {
	awsCfg := aws.Config{
		Region:      cfg.S3.Region,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.S3.AccessKeyID, cfg.S3.SecretAccessKey, ""),
	}

	opts := []func(*s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.S3.UsePathStyle },
	}
	if cfg.S3.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, opts...)
	return &S3Client{
		client:    client,
		presigner: s3.NewPresignClient(client),
		bucket:    cfg.S3.Bucket,
		cdnPrefix: cfg.CDN,
	}
}

// PresignPut returns a presigned PUT URL for direct browser upload.
func (s *S3Client) PresignPut(ctx context.Context, key string) (string, error) {
	req, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}
	return req.URL, nil
}

// PublicURL returns the CDN or S3 URL for a key.
func (s *S3Client) PublicURL(key string) string {
	if s.cdnPrefix != "" {
		return s.cdnPrefix + "/" + key
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
}

// Bucket returns the configured bucket name.
func (s *S3Client) Bucket() string { return s.bucket }

// DeleteObject deletes a key from the given bucket.
func (s *S3Client) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
