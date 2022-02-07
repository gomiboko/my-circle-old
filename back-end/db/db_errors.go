package db

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

type ErrorNumber uint16

const ErrDuplicateEntry ErrorNumber = 1062

func Is(err error, errNum ErrorNumber) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == uint16(errNum)
}
