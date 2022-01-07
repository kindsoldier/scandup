package pmscanner

import (
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"

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

const maxDepth int = 7

func (this *Scanner) Scan(basedir string) error {
    var err error

    err = this.keyvdb.Clean()
	if err != nil {
		return err
	}

    executor := func(filename string) {
        fmt.Print(".")
        this.scanFile(filename)
    }
    _ = this.FileTreeScanner(basedir, maxDepth, executor)
    fmt.Println()
    return err
}

func (this *Scanner) scanFile(filename string) error {
    var err error

    if this.keyvdb == nil {
        return errors.New("db yet not open")
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
        if  err != nil {
            return err
        }
        //if read < this.bSize {
        //    continue
        //}
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
    dupSize := 0
    dupCounter := 0
    executor := func(key []byte, val []byte) bool {
        counter, _ := strconv.Atoi(string(val))
        if counter > 1 {
            dupSize += (counter - 1) * this.bSize
            dupCounter += (counter - 1)
        }
        return false
    }
    err = this.keyvdb.Iter(executor)
    if err != nil {
        return err
    }

    totalCounter := 0
    executor = func(key []byte, val []byte) bool {
        totalCounter++
        return false
    }
    err = this.keyvdb.Iter(executor)
    if err != nil {
        return err
    }
    fmt.Println(totalCounter, dupCounter, (100 * float32(dupCounter))/float32(totalCounter))
    //totalSize := totalCounter * this.bSize
    //fmt.Println(totalSize, dupSize, (100 * float32(dupSize))/float32(totalSize))

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

const pathSeparator string = "/"

type FileExecutor = func(filePath string)

func (this *Scanner) FileTreeScanner(basePath string, depth int, executor FileExecutor) error {

    pathLength := func(path string) int {
        if len(path) == 0 {
            return 0
        }
        path = filepath.Clean(path)
        path = filepath.ToSlash(path)
        path = strings.Trim(path, pathSeparator)
        if len(path) == 0 {
            return 0
        }
        list := strings.Split(path, pathSeparator)
        return len(list)
    }

    depth = depth + pathLength(basePath)

    resolver := func(fullPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if pathLength(fullPath) > depth  {
            return filepath.SkipDir
        }
        if !info.IsDir(){
            executor(fullPath)
        }
        return nil
    }
    err := filepath.Walk(basePath, resolver)
    if err != nil {
        return err
    }
    return err
}



//EOF
