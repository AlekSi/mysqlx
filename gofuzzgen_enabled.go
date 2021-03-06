// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017-2018 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gofuzz

package mysqlx

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func addToGoFuzzCorpus(name string, data []byte) {
	path := filepath.Join("go-fuzz", name, "corpus")
	_ = os.MkdirAll(path, 0777)

	path = filepath.Join(path, fmt.Sprintf("test-%x", sha1.Sum(data)))
	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		panic(err)
	}
}
