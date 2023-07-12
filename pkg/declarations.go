package pkg

import (
	"database/sql"

	"github.com/rs/zerolog"
)

// We're creating a single struct to persist a single connection to the db
// instead of opening a new one every time we run a query. This is more efficient
// and should take up less time for each query. In this smaller project, we
// wont see any notable differences, but at a larger scale, this will matter.
type Conn struct {
	Db   *sql.DB
	Zlog zerolog.Logger
}

type Ship struct {
	Id     int
	Name   string  `json:"name"`
	Class  string  `json:"class"`
	Crew   int     `json:"crew"`
	Img    string  `json:"img"`
	Value  float32 `json:"value"`
	Status string  `json:"status"`
	Arms   []*Armament
}

type Armament struct {
	Name string
	Qty  int
}
