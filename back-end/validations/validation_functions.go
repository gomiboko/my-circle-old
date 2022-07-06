package validations

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
)

const maxIconFileSize = 1_000_000

func IsOverMaxIconFileSize(fh *multipart.FileHeader) bool {
	return fh != nil && fh.Size > maxIconFileSize
}

func IsNotAllowedIconFileFormat(fh *multipart.FileHeader) bool {
	file, err := fh.Open()
	if err != nil {
		log.Print("ファイルオープンエラー")
		return true
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	// importしていない形式の画像ファイル、または画像以外のファイルの場合
	if err != nil {
		return true
	}

	return format != "png" && format != "jpeg"
}
