package model

import (
	"flag"
)

var (
	// Host of database server
	Host = flag.String(`dbHost`, ``, `Database Host`)

	// Port of database server
	Port = flag.String(`dbPort`, `5432`, `Database Port`)

	// User of database
	User = flag.String(`dbUser`, `funds`, `Database User`)

	// Pass of database
	Pass = flag.String(`dbPass`, ``, `Database Pass`)

	// Name of database
	Name = flag.String(`dbName`, `funds`, `Database Name`)
)
