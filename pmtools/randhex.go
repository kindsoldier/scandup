//
package pmtools

import (
    "encoding/hex"
    "time"
    "math/rand"
)

func RandBytesHex(size int) string {
    rand.Seed(time.Now().UnixNano())
    randBytes := make([]byte, size)
    rand.Read(randBytes)
    hexString := hex.EncodeToString(randBytes)
    return hexString
}

func RandBytes(size int) []byte {
    rand.Seed(time.Now().UnixNano())
    randBytes := make([]byte, size)
    rand.Read(randBytes)
    return randBytes
}
//EOF
