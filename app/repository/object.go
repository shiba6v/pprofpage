package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/shiba6v/eu"
)

type ObjectStorage struct {
	s3cli      *s3.Client
	bucketName string
}

func NewObjectStorage(s3cli *s3.Client, bucketName string) ObjectStorage {
	return ObjectStorage{s3cli: s3cli, bucketName: bucketName}
}

func CleansePath(s string) string {
	// directory traversal
	return strings.Trim(s, "/")
}

func (s ObjectStorage) GetObjectToTmp(ctx context.Context, key string) (string, error) {
	key = CleansePath(key)
	obj, err := s.s3cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", eu.Wrap(err)
	}
	path := fmt.Sprintf("/tmp/%s", key)
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		// use cache
		return path, nil
	}
	newFile, err := os.Create(path)
	if err != nil {
		return "", eu.Wrap(err)
	}
	defer newFile.Close()
	if _, err := io.Copy(newFile, obj.Body); err != nil {
		return "", eu.Wrap(err)
	}
	return path, nil
}

func (s ObjectStorage) UploadObject(ctx context.Context, key string, b []byte) error {
	key = CleansePath(key)
	path := fmt.Sprintf("/tmp/%s", key)
	newFile, err := os.Create(path)
	if err != nil {
		return eu.Wrap(err)
	}
	defer newFile.Close()
	if _, err := io.Copy(newFile, bytes.NewReader(b)); err != nil {
		return eu.Wrap(err)
	}
	_, err = s.s3cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return eu.Wrap(err)
	}
	return nil
}
