package videos

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(ctx context.Context) *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("ap-northeast-1"),
	)
	if err != nil {
		panic(err)
	}
	return s3.NewFromConfig(cfg)
}

func Uploadfile(ctx context.Context, client *s3.Client, bucket string) {

}

func Downloadfile(ctx context.Context, client *s3.Client, key string) *bytes.Reader {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("elasticbeanstalk-ap-northeast-1-271366295554"),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("error:%v", err)
	}
	defer out.Body.Close()

	buff, err := io.ReadAll(out.Body)
	if err != nil {
		log.Printf("error:%v", err)
	}
	reader := bytes.NewReader(buff)
	return reader
}
