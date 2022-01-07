//
package main

import (
    "fmt"
    "pmapp/pmscanner"
)

func main() {
    err := run()
    if err != nil {
        fmt.Println("error:", err)
    }
}

const bSize int     = 512
var key     []byte  = []byte("01234567890012345678900123456789012")

func run() error {
    var err error
    scanner, err := pmscanner.NewScanner("tmpdb", key, bSize)
    if err != nil {
        return err
    }
    err = scanner.Scan("data.bin")
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




