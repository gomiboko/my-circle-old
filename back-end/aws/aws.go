package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gomiboko/my-circle/consts"
)

var conf aws.Config

func Init() error {
	localStackEndpoint := os.Getenv(consts.SessKeyLocalStackEndpoint)
	awsRegion := os.Getenv(consts.SessKeyAWSRegion)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if localStackEndpoint != "" {
			// 開発環境用のエンドポイント情報を返却
			return aws.Endpoint{
				PartitionID:   consts.AwsPartitionID,
				URL:           localStackEndpoint,
				SigningRegion: awsRegion,
			}, nil
		} else {
			// EndpointNotFoundErrorを返却することで、デフォルトのリゾルバを使用する
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		}
	})

	var err error
	conf, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	return err
}

func GetConf() aws.Config {
	return conf
}
