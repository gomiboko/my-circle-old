package utils

const condRowVersion = "row_version = ?"

func CreateRowVersionCond(rowVersion uint) (string, uint) {
	return condRowVersion, rowVersion
}
