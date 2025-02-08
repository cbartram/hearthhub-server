package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
)

type S3Service struct {
	client *s3.Client
	bucket string
}

// MakeS3Service creates a new instance of S3Service
func MakeS3Service(region string) (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Service{client: client, bucket: os.Getenv("BUCKET")}, nil
}

// ListObjects lists all objects in a bucket with given prefix
func (s *S3Service) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	var objects []string
	paginator := s3.NewListObjectsV2Paginator(s.client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}

		for _, obj := range page.Contents {
			objects = append(objects, *obj.Key)
		}
	}

	return objects, nil
}

// UploadObject publishes an object to S3 under the given prefix
func (s *S3Service) UploadObject(ctx context.Context, key string) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put object: %v", err)
	}

	return result, nil
}

// DeleteObject deletes an object from S3
func (s *S3Service) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

// DeleteObjectsWithPrefix deletes all objects with the given prefix
func (s *S3Service) DeleteObjectsWithPrefix(ctx context.Context, prefix string) error {
	// First list all objects with the prefix
	objects, err := s.ListObjects(ctx, prefix)
	if err != nil {
		return fmt.Errorf("failed to list objects for deletion: %v", err)
	}

	// Create delete objects input
	var objectIds []types.ObjectIdentifier
	for _, key := range objects {
		objectIds = append(objectIds, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	// Delete objects in batches of 100
	const maxBatchSize = 100
	for i := 0; i < len(objectIds); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(objectIds) {
			end = len(objectIds)
		}

		batch := objectIds[i:end]
		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &types.Delete{
				Objects: batch,
				Quiet:   aws.Bool(true),
			},
		}

		_, err := s.client.DeleteObjects(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to delete objects batch: %v", err)
		}
	}

	return nil
}
