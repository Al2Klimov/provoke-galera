package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
)

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
}
