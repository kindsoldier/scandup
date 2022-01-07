//
package pmkeyvdb

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
    db *leveldb.DB
}

func NewDB() *DB {
    return &DB{}
}

func (this *DB) Open(filename string) error {
    var err error
    if this.db != nil {
        this.Close()
    }
    db, err := leveldb.OpenFile(filename, nil)
    if err != nil {
         return errors.New(fmt.Sprintf("cannot open database: %v", err))
    }
    this.db = db
    return err
}

func (this *DB) Set(key []byte, value []byte) error {
    var err error
    err = this.db.Put(key, value, nil)
    if err != nil {
        return errors.New("db not yet open")

        return err
    }
    return err
}

func (this *DB) Get(key []byte) ([]byte, error) {
    var err   error
    var value []byte
    value, err = this.db.Get(key, nil)
    if err != nil {
        return value, err
    }
    return value, err
}

func (this *DB) Has(key []byte) (bool, error) {
    var err  error
    var has bool
    if this.db == nil {
        return has, errors.New("db not yet open")
    }
    has, err = this.db.Has(key, nil)
    if err != nil {
        return has, err
    }
    return has, err
}

type Resolver = func(key []byte, val []byte) (ok bool)

func (this *DB) First(sval []byte) ([]byte, bool, error) {
    compf := func(key []byte, val []byte) (ok bool) {
        if bytes.Equal(val, sval) {
            return true
        }
        return false
    }
    return this.comp(compf)
}

func (this *DB) comp(resolver Resolver) ([]byte, bool, error) {
    var err  error
    var key  []byte
    var ok bool

    iter := this.db.NewIterator(nil, nil)
    defer iter.Release()
    for iter.Next() {
        if resolver(iter.Key(), iter.Value()) {
            ok = true
            key = iter.Key()
            break
        }
    }
    err = iter.Error()
    if err != nil {
        return key, ok, err
    }
    return key, ok, err
}


func (this *DB) All(resolver Resolver) ([][]byte, bool, error) {
    collect := make([][]byte, 0)
    var ok bool
    var err error

    iter := this.db.NewIterator(nil, nil)
    defer iter.Release()
    for iter.Next() {
        if resolver(iter.Key(), iter.Value()) {
            ok = true
            collect = append(collect, iter.Key())
        }
    }
    err = iter.Error()
    if err != nil {
        return collect, ok, err
    }
    return collect, ok, err
}

type Executor = func(key []byte, val []byte) (stop bool)

func (this *DB) Iter(executor Executor) error {
    var err error
    iter := this.db.NewIterator(nil, nil)
    defer iter.Release()
    for iter.Next() {
        if executor(iter.Key(), iter.Value()) {
            break
        }
    }
    err = iter.Error()
    if err != nil {
        return err
    }
    return err
}

func (this *DB) Clean() error {
    var err error
    iter := this.db.NewIterator(nil, nil)
    defer iter.Release()
    for iter.Next() {
        err = this.db.Delete(iter.Key(), nil)
        if err != nil {
            return err
        }
    }
    err = iter.Error()
    if err != nil {
        return err
    }
    return err
}


func (this *DB) Close() {
    if this.db != nil {
        this.db.Close()
    }
}
//EOF
