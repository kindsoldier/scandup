//
package main

import (
    "fmt"
    "path/filepath"
    "os"
    "pmapp/pmscanner"
)

func main() {

    err := run()
    if err != nil {
        fmt.Println("error:", err)
    }
}

const bSize  int     = 1024 * 4
const dbname string  = "tmp.db"
var key      []byte  = []byte("01234567890012345678900123456789012")

func help() {
    fmt.Println("usage:", filepath.Base(os.Args[0]), "dbdir scandir")
}

func run() error {
    var err error

    if len(os.Args) < 3 {
        help()
        return err
    }
    dbdir := os.Args[1]
    scandir := os.Args[2]

    dbpath := filepath.Join(dbdir, dbname)
    scanner, err := pmscanner.NewScanner(dbpath, key, bSize)
    if err != nil {
        return err
    }
    err = scanner.Scan(scandir)
    if err != nil {
        return err
    }
    err = scanner.Print()
    if err != nil {
        return err
    }
    return err
}
//EOF




