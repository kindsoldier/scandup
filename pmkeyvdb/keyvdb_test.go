//
package pmkeyvdb

import (
    "testing"
    "path/filepath"

    "pmapp/tools"
    "github.com/stretchr/testify/assert"
)

const dbname string = "tmp.leveldb"
const repeat int    = 1024

func TestHasher01SetGetRandKey(t *testing.T) {
    var err error
    dbpath := filepath.Join(t.TempDir(), dbname)
    db := NewDB()
    err = db.Open(dbpath)
    defer db.Close()
    if err != nil {
        t.Fatal(err)
    }
    for i := 0; i < repeat; i++ {
        key := tools.RandBytes(16)
        ival := tools.RandBytes(128)
        err = db.Set(key, ival)
        if err != nil {
            t.Fatal(err)
        }
        oval, err := db.Get(key)
        if err != nil {
            t.Fatal(err)
        }
        assert.Equal(t, ival, oval, nil)
    }
}

func TestHasher02SetGetEqKey(t *testing.T) {
    var err error
    dbpath := filepath.Join(t.TempDir(), dbname)
    db := NewDB()
    err = db.Open(dbpath)
    defer db.Close()
    if err != nil {
        t.Fatal(err)
    }
    key := tools.RandBytes(128)
    for i := 0; i < repeat; i++ {
        ival := tools.RandBytes(128)
        err = db.Set(key, ival)
        if err != nil {
            t.Fatal(err)
        }
        oval, err := db.Get(key)
        if err != nil {
            t.Fatal(err)
        }
        assert.Equal(t, ival, oval, nil)
    }
}
//EOF

