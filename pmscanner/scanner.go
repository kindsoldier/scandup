package pmscanner

import (
    "errors"
    "fmt"
    "os"
    "io"
    "strconv"
    "encoding/hex"

    "pmapp/pmhasher"
    "pmapp/pmkeyvdb"
)


type Scanner struct {
    bSize   int
    keyvdb  *pmkeyvdb.DB
    hasher  *pmhasher.Hasher
}

func NewScanner(dbname string, key []byte, bSize int) (*Scanner, error) {
    var err error
    scanner := &Scanner{}

    keyvdb := pmkeyvdb.NewDB()
    err = keyvdb.Open(dbname)
    if err != nil {
        return scanner, err
    }
    scanner.keyvdb = keyvdb

    hasher, err := pmhasher.NewHahser(key)
    if err != nil {
        return scanner, err
    }
    scanner.hasher = hasher
    scanner.bSize = bSize
    return scanner, err
}

func (this *Scanner) Close() {
    this.keyvdb.Close()
}

func (this *Scanner) Scan(filename string) error {
    var err error

    if this.keyvdb == nil {
        return errors.New("db yet not open")
    }

    err = this.keyvdb.Clean()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

    for {
        buffer := make([]byte, this.bSize)
        read, err := file.Read(buffer);
        if  err == io.EOF {
            break
        }
        hash, err := this.hasher.Hash(buffer[:read])
        if err != nil {
            return err
        }
        err = this.Inc(hash)
        if err != nil {
            return err
        }
    }
    return err
}

func (this *Scanner) Print() error {
    var err error
    executor := func(key []byte, val []byte) bool {
        keyHex := hex.EncodeToString(key)
        counter, _ := strconv.Atoi(string(val))
        if counter > 1 {
            fmt.Println(keyHex, counter)
        }
        return false
    }
    err = this.keyvdb.Iter(executor)
    if err != nil {
        return err
    }
    return err
}

func (this *Scanner) Inc(key []byte) error {
    var err error
    var val []byte
    has, err := this.keyvdb.Has(key)
    if err != nil {
        return err
    }
    if has {
        val, err = this.keyvdb.Get(key)
        if err != nil {
            return err
        }
    }
    var counter int
    switch {
        case len(val) == 0:
            counter = 1
        default:
            counter, err = strconv.Atoi(string(val))
            if err != nil {
                return err
            }
            counter++
    }
    nval := []byte(strconv.Itoa(counter))
    err = this.keyvdb.Set(key, nval)
    if err != nil {
        return err
    }
    return err
}
//EOF
