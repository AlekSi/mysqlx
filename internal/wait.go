// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/AlekSi/mysqlx"
)

const timeout = 60 * time.Second

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s waits for MySQL to become available.\n", os.Args[0])
	}
	flag.Parse()
	log.SetFlags(log.Lmicroseconds)

	dataSource := os.Getenv("MYSQLX_TEST_DATASOURCE")
	if dataSource == "" {
		log.Fatal("Please set environment variable MYSQLX_TEST_DATASOURCE.")
	}

	log.Printf("Connecting to %s ...", dataSource)
	start := time.Now()
	var prevErr error
	var attempts int
	for {
		attempts++
		db, err := sql.Open("mysqlx", dataSource)
		if err == nil {
			err = db.Ping()
		}

		if err != nil {
			if prevErr != err || attempts%10 == 0 {
				log.Print(err)
				prevErr = err
			}
			if time.Since(start) > timeout {
				log.Fatalf("Failed! Last error: %s", err)
			}
			time.Sleep(time.Second)
			continue
		}

		log.Print("Connected!")
		db.Close()
		return
	}
}
