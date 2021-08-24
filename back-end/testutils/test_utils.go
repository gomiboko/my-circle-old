package testutils

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetFixtures(fixturesDirPath string) (*testfixtures.Loader, error) {
	sqldb, err := sql.Open("mysql", getDSN())
	if err != nil {
		return nil, err
	}

	return testfixtures.New(
		testfixtures.Database(sqldb),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(fixturesDirPath),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
}

func GetDB() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(getDSN()), &gorm.Config{})
}

func getDSN() string {
	return fmt.Sprintf("root:root@tcp(test-db:3306)/mycircle?charset=utf8mb3&parseTime=True&loc=%s",
		url.QueryEscape("Asia/Tokyo"))
}
