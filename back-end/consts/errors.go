package consts

import "errors"

var (
	ErrS3KeyRequired   = errors.New("S3へのアップロードにはキー値の指定が必要です")
	ErrS3FileRequired  = errors.New("S3へのアップロードにはファイルの指定が必要です")
	ErrUpdatedByOthers = errors.New("表示されていたデータが古い為、更新できませんでした")
)
