package services

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const testFilePath = "./testdata/no1.png"

type StorageServiceTestSuite struct {
	suite.Suite
	storageService StorageService
	s3Client       *s3.Client
}

func (s *StorageServiceTestSuite) SetupSuite() {
	conf, err := testutils.GetAwsConfig()
	if err != nil {
		s.FailNow(err.Error())
	}

	s.storageService = NewS3Service(conf)

	// テスト結果確認用のS3クライアント
	s.s3Client = s3.NewFromConfig(conf, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func TestStorageService(t *testing.T) {
	suite.Run(t, new(StorageServiceTestSuite))
}

func (s *StorageServiceTestSuite) TestUpload() {
	s.Run("不正な入力値の場合", func() {
		fh, err := testutils.CreateFileHeader(testFilePath)
		require.Nil(s.T(), err)

		inputs := []struct {
			key      string
			fh       *multipart.FileHeader
			expected error
		}{
			{key: "", fh: fh, expected: consts.ErrS3KeyRequired},
			{key: "key", fh: nil, expected: consts.ErrS3FileRequired},
		}

		for _, in := range inputs {
			err = s.storageService.Upload(in.key, in.fh)
			assert.True(s.T(), errors.Is(err, in.expected))
		}
	})

	s.Run("ファイルのオープンに失敗した場合", func() {
		wrongFh := &multipart.FileHeader{}
		err := s.storageService.Upload("key", wrongFh)

		require.NotNil(s.T(), err)
		require.Contains(s.T(), err.Error(), consts.ErrMsgFailedToLoadFile)
	})

	s.Run("ファイルのアップロードに失敗した場合", func() {
		fh, err := testutils.CreateFileHeader(testFilePath)
		require.Nil(s.T(), err)

		// 不正な aws.Config で StorageService を作成してアップロード
		wrongStorageService := NewS3Service(aws.Config{})
		err = wrongStorageService.Upload("key", fh)

		require.NotNil(s.T(), err)
		require.Contains(s.T(), err.Error(), consts.ErrMsgFailedToRegisterToS3)
	})

	s.Run("ファイルのアップロードに成功した場合", func() {
		fh, err := testutils.CreateFileHeader(testFilePath)
		require.Nil(s.T(), err)

		const key = "key"
		err = s.storageService.Upload(key, fh)

		require.Nil(s.T(), err)

		// LocalStackに登録されていることを確認
		_, err = s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(consts.AwsS3BucketName),
			Key:    aws.String(key),
		})
		require.Nil(s.T(), err)
	})
}
