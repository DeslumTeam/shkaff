package psql

// import (
// 	"fmt"
// 	"log"

// 	"github.com/jmoiron/sqlx"
// )

// const (
// 	URI_TEMPLATE = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
// )

// type pSQL struct {
// 	uri             string
// 	db              *sqlx.DB
// 	refreshTimeScan int
// }

// func initPSQL(user, password, host string, port int, database_db string) (ps *pSQL) {
// 	var err error
// 	ps = new(pSQL)
// 	ps.uri = fmt.Sprintf(URI_TEMPLATE, user, password, host, port, database_db)
// 	ps.refreshTimeScan = 10
// 	if ps.db, err = sqlx.Connect("postgres", ps.uri); err != nil {
// 		log.Fatalln(err)
// 	}
// 	return
// }
