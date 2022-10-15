package utils

const condRowVersion = "row_version = ?"

func CreateRowVersionCond(rowVersion uint) (string, uint) {
	return condRowVersion, rowVersion
}

func CreateOrderStr(columns ...string) string {
	columnsStr := ""
	for i, col := range columns {
		if i == 0 {
			columnsStr += col
		} else {
			columnsStr += (", " + col)
		}
	}

	return columnsStr
}
