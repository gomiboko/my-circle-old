package validations

import "mime/multipart"

const maxIconFileSize = 1_000_000

func IsOverMaxIconFileSize(fh *multipart.FileHeader) bool {
	return fh != nil && fh.Size > maxIconFileSize
}
