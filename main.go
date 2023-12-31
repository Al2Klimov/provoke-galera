package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"net/url"
	"strings"
)

//go:embed schema.sql
var schema string

func main() {
	host1 := flag.String("host1", "", "HOST")
	host2 := flag.String("host2", "", "HOST")
	db := flag.String("db", "", "DATABASE")
	user := flag.String("user", "", "USERNAME")
	pass := flag.String("pass", "", "PASSWORD")
	flag.Parse()

	if *host1 == "" {
		panic("-host1 missing")
	}

	if *host2 == "" {
		panic("-host2 missing")
	}

	if *db == "" {
		panic("-db missing")
	}

	if *user == "" {
		panic("-user missing")
	}

	if *pass == "" {
		panic("-pass missing")
	}

	db1, err1 := sql.Open("mysql", (&url.URL{
		User: url.UserPassword(*user, *pass),
		Host: "tcp(" + net.JoinHostPort(*host1, "3306") + ")",
		Path: "/" + *db,
	}).String()[2:])
	if err1 != nil {
		panic(err1)
	}
	defer func() { _ = db1.Close() }()

	db2, err2 := sql.Open("mysql", (&url.URL{
		User: url.UserPassword(*user, *pass),
		Host: "tcp(" + net.JoinHostPort(*host2, "3306") + ")",
		Path: "/" + *db,
		//RawQuery: "wsrep_sync_wait=4",
	}).String()[2:])
	if err2 != nil {
		panic(err2)
	}
	defer func() { _ = db2.Close() }()

	for _, ddl := range strings.Split(strings.Trim(strings.TrimSpace(schema), ";"), ";") {
		if _, err := db1.Exec(ddl); err != nil {
			panic(err)
		}
	}

	parent, errP := db1.Prepare("INSERT INTO `parent` VALUES ()")
	if errP != nil {
		panic(errP)
	}
	defer func() { _ = parent.Close() }()

	child, errC := db2.Prepare("INSERT INTO `child`(`id`) VALUES (?)")
	if errC != nil {
		panic(errC)
	}
	defer func() { _ = child.Close() }()

	for {
		result, errP := parent.Exec()
		if errP != nil {
			panic(errP)
		}

		id, errLI := result.LastInsertId()
		if errLI != nil {
			panic(errLI)
		}

		for {
			if _, err := child.Exec(id); err != nil {
				if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1213 {
					continue
				}

				panic(err)
			}

			break
		}

		println(id)
	}
}
