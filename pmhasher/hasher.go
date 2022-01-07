package pmhasher

import (
    "errors"
    "github.com/minio/highwayhash"
)


type Hasher struct {
    key []byte
}

const keySize int = 32

func NewHahser(key []byte) (*Hasher, error) {
    var err error
    hasher := &Hasher{}
    if len(key) < keySize {
        return hasher, errors.New("too short init key")
    }
    hasher.key = key[0:keySize]
    return hasher, err
}

func (this *Hasher) Hash(data []byte) ([]byte, error) {
    var err error
    var sum []byte
    hash, err := highwayhash.New(this.key)
    if err != nil {
        return sum, err
    }
    hash.Write(data)
    sum = hash.Sum(nil)
    return sum, err
}
//EOF
