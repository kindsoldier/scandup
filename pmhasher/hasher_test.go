//
package pmhasher


import (
    "encoding/hex"
    "testing"

    "pmapp/tools"
    "github.com/stretchr/testify/assert"
)

const hashHexSample string  = "5c5278baa03952257df4db5dcccc2e861dcf211094aac9ae1a3c7fdbc10fe0d3"
const repeat1        int     = 1024 * 64
const repeat2        int     = 1024 * 64

func TestHasher01Pre(t *testing.T) {
    key := []byte("1234567890123456789012345678901234567890")
    hasher, err := NewHahser(key)
    if err != nil {
        t.Fatal(err)
    }
    for i := 0; i < repeat1; i++ {
        data := []byte("1234567890123456789012345678901234567890")
        hash, err := hasher.Hash(data)
        if err != nil {
            t.Fatal(err)
        }
        hashHex := hex.EncodeToString(hash)
        assert.Equal(t, hashHexSample, hashHex, nil)
    }
}

func TestHasher02Equal(t *testing.T) {
    key := tools.RandBytes(32)
    hasher, err := NewHahser(key)
    if err != nil {
        t.Fatal(err)
    }
    for i := 0; i < repeat1; i++ {
        data := tools.RandBytes(1024)
        hash1, err := hasher.Hash(data)
        if err != nil {
            t.Fatal(err)
        }
        hash2, err := hasher.Hash(data)
        if err != nil {
            t.Fatal(err)
        }
        hash1Hex := hex.EncodeToString(hash1)
        hash2Hex := hex.EncodeToString(hash2)
        assert.Equal(t, hash1Hex, hash2Hex, nil)
    }
}

func TestHasher03Rand(t *testing.T) {
    key := tools.RandBytes(32)
    hasher, err := NewHahser(key)
    if err != nil {
        t.Fatal(err)
    }
    data1 := tools.RandBytes(16)
    hash1, err := hasher.Hash(data1)
    if err != nil {
        t.Fatal(err)
    }
    hash1Hex := hex.EncodeToString(hash1)

    for i := 0; i < repeat2; i++ {
        data2 := tools.RandBytes(16)
        hash2, err := hasher.Hash(data2)
        if err != nil {
            t.Fatal(err)
        }
        hash2Hex := hex.EncodeToString(hash2)
        assert.NotEqual(t, hash1Hex, hash2Hex, nil)
    }
}
//EOF
