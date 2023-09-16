package registry

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shiba6v/pprofpage/app/controller"
	"github.com/shiba6v/pprofpage/app/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func RegisterServer(e *echo.Echo, s3cli *s3.Client) error {
	bucketName := Getenv("BUCKET_NAME")
	storage := repository.NewObjectStorage(s3cli, bucketName)
	r := controller.NewController(storage)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status} ${error}\n",
	}))
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/pprof/:id", r.GetPProf)
	e.GET("/pprof/:id/:key", r.GetPProf)
	e.POST("/pprof/register", r.RegisterPProf)
	e.GET("/health", r.Health)
	return nil
}

func Getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(fmt.Sprintf("os.Getenv(%s) not found", name))
	}
	return v
}

func RegisterMinio() (*s3.Client, error) {
	ctx := context.Background()
	accessKey := Getenv("MINIO_ACCESS_KEY")
	secretKey := Getenv("MINIO_SECRET_KEY")
	cred := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	endpoint := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "http://minio:9100",
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cred), config.WithEndpointResolverWithOptions(endpoint))
	if err != nil {
		log.Fatalln(err)
	}
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	return client, nil
}

func RegisterS3() (*s3.Client, error) {
	ctx := context.Background()
	endpoint := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	// https://zenn.dev/papu_nika/articles/b8efdef1c87f65
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"), config.WithEndpointResolverWithOptions(endpoint))
	if err != nil {
		log.Fatalln(err)
	}
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
	})
	return client, nil
}
