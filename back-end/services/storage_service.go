package services

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/utils"
)

type StorageService interface {
	Upload(key string, fh *multipart.FileHeader) error
}

type s3Service struct {
	conf aws.Config
}

func NewS3Service(conf aws.Config) StorageService {
	return &s3Service{
		conf: conf,
	}
}

func (s3s *s3Service) Upload(key string, fh *multipart.FileHeader) error {
	if key == "" {
		return consts.ErrS3KeyRequired
	}

	if fh == nil {
		return consts.ErrS3FileRequired
	}

	file, err := fh.Open()
	if err != nil {
		return utils.NewErrorWithInnerError(consts.ErrMsgFailedToLoadFile, err)
	}
	defer file.Close()

	client := s3.NewFromConfig(s3s.conf, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(consts.AwsS3BucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return utils.NewErrorWithInnerError(consts.ErrMsgFailedToRegisterToS3, err)
	}

	return nil
}
