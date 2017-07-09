// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/AlekSi/mysqlx"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s waits for MySQL to become available.\n", os.Args[0])
	}
	flag.Parse()
	log.SetFlags(log.Lmicroseconds)

	ds := os.Getenv("MYSQLX_TEST_DATASOURCE")
	if ds == "" {
		log.Fatal("Please set environment variable MYSQLX_TEST_DATASOURCE.")
	}
	u, err := url.Parse(ds)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to %s ...", u)
	start := time.Now()
	for {
		db, err := sql.Open("mysqlx", u.String())
		if err == nil {
			err = db.Ping()
		}

		if err != nil {
			log.Print(err)
			if time.Since(start) > 30*time.Second {
				log.Fatal("Failed!")
			}
			time.Sleep(2 * time.Second)
			continue
		}

		log.Print("Connected!")
		db.Close()
		return
	}
}
