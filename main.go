package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	_ "github.com/alexbrainman/odbc"
)

var (
	mssrv    = flag.String("mssrv", "server", "ms sql server name")
	msdb     = flag.String("msdb", "dbname", "ms sql server database name")
	msuser   = flag.String("msuser", "", "ms sql server user name")
	mspass   = flag.String("mspass", "", "ms sql server password")
	msdriver = flag.String("msdriver", "sql server", "ms sql odbc driver name")
	msport   = flag.String("msport", "1433", "ms sql server port number")
)

type connParams map[string]string

func newConnParams() connParams {
	params := connParams{
		"driver":   *msdriver,
		"server":   *mssrv,
		"database": *msdb,
	}

	if len(*msuser) == 0 {
		params["trusted_connection"] = "yes"
	} else {
		params["uid"] = *msuser
		params["pwd"] = *mspass
	}

	a := strings.SplitN(params["server"], ",", -1)
	if len(a) == 2 {
		params["server"] = a[0]
		params["port"] = a[1]
	}
	return params
}

func mssqlConnectWithParams(params connParams) (db *sql.DB, stmtCount int, err error) {
	db, err = sql.Open("odbc", params.makeODBCConnectionString())
	if err != nil {
		return nil, 0, err
	}
	stats := db.Driver().(*Driver).Stats
	return db, stats.StmtCount, nil
}

func (params connParams) makeODBCConnectionString() string {
	if port, ok := params["port"]; ok {
		params["server"] += "," + port
		delete(params, "port")
	}
	var c string
	for n, v := range params {
		c += n + "=" + v + ";"
	}
	return c
}

func main() {
	//db, err = sql.open("odbc", )
	flag.Parse()
	fmt.Printf("servername: %s, dbname: %s, user: %s, password: %s, msdriver: %s, msport: %s\n",
		*mssrv,
		*msdb,
		*msuser,
		*mspass,
		*msdriver,
		*msport)

	params := newConnParams()

	db, sc, err := mssqlConnectWithParams(params)
}
