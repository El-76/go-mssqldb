package mssql

import (
	"github.com/golang-sql/civil"
)

type NullDateTime struct {
	DateTime civil.DateTime
	Valid    bool
}

type NullDate struct {
	Date  civil.Date
	Valid bool
}

type NullTime struct {
	Time  civil.Time
	Valid bool
}
