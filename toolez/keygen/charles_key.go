package keygen

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type request struct {
	name string `binding: required`
}

// GetCharlesKey process http Get request
func GetCharlesKey(c *gin.Context) {
	var req request
	if c.Bind(&req) == nil {
		key := CharlesKeygen(req.name)
		c.JSON(http.StatusOK, gin.H{"key": key})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not found"})
	}
}

var a = make([]uint32, 26)
var b = make([]uint32, 26)

// CharlesKeygen return a string as key for the name
// name -> x
// [x, 01, random]=s=n -> crc
// [x, 01, random] -> y=m
// key: 'crc' + 'y'
func CharlesKeygen(name string) string {
	load(8800536498351690864)
	x := encryptName(name) ^ 0x54882F8A
	r := rand.Int31n(0xFFFFF) + 0x100000
	s := fmt.Sprintf("%x01%x", x, r)
	n, _ := strconv.ParseUint(s, 16, 64)
	c := fmt.Sprintf("%02x", genCRC(n))
	load(13038168091975921581)
	y := unencryptUint64(n)
	m := fmt.Sprintf("%x", y)
	return fmt.Sprintf("%s%s", c, m)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	var x int32
	x = -1209970333
	a[0] = uint32(x)

	x = -1209970333
	b[0] = uint32(x)
	for i := 1; i < 26; i++ {
		x -= 1640531527
		b[i] = uint32(x)
	}
}

func genCRC(n uint64) uint32 {
	var t uint64
	for i := 56; i >= 0; i -= 8 {
		t ^= n >> uint(i) & 255
	}
	return uint32(t & 255)
}

/************* encrypt name *****************/

func encryptName(name string) uint32 {
	lenOld := len(name)
	lenNew := lenOld + 4
	if lenNew%8 != 0 {
		lenNew += 8 - lenNew%8
	}

	newName := make([]byte, lenNew)
	copy(newName[4:], name)
	newName[0] = byte(lenOld >> 24)
	newName[1] = byte(lenOld >> 16)
	newName[2] = byte(lenOld >> 8)
	newName[3] = byte(lenOld)

	newName = encryptArray(newName)

	var ret uint32
	for i := 0; i < lenNew; i++ {
		ret ^= uint32(int8(newName[i]))
		ret = ret<<3 | ret>>29
	}

	return ret
}

// combine every 4 bytes into uint64, encrypt it,
// and unpack back to 4 bytes
func encryptArray(newName []byte) []byte {
	// mixed name
	mn := make([]byte, len(newName))
	var t uint64

	for i := 1; i <= len(newName); i++ {
		t = t<<8 | uint64(newName[i-1])
		if i%8 == 0 {
			t = encryptUint64(t)
			mn[i-8] = byte(t >> 56)
			mn[i-7] = byte(t >> 48)
			mn[i-6] = byte(t >> 40)
			mn[i-5] = byte(t >> 32)
			mn[i-4] = byte(t >> 24)
			mn[i-3] = byte(t >> 16)
			mn[i-2] = byte(t >> 8)
			mn[i-1] = byte(t)
			t = 0
		}
	}

	return mn
}

func encryptUint64(n uint64) uint64 {
	lo := uint32(n) + a[0]
	hi := uint32(n>>32) + a[1]
	var t, k uint32

	var f = func(i int, x, y *uint32) {
		t = *x ^ *y
		k = *x & 31
		*y = (t<<k | t>>(32-k)) + a[i+2]
	}

	for i := 0; i < 24; i++ {
		if i%2 == 0 {
			f(i, &hi, &lo)
		} else {
			f(i, &lo, &hi)
		}
	}

	return uint64(hi)<<32 + uint64(lo)
}

func unencryptUint64(n uint64) uint64 {
	lo := uint32(n)
	hi := uint32(n >> 32)
	var t, k uint32

	var f = func(i int, x, y *uint32) {
		k = *x & 31
		t = *y - a[25-i]
		xor := t<<(32-k) | t>>k
		*y = xor ^ *x
	}
	for i := 0; i < 24; i++ {
		if i%2 == 0 {
			f(i, &lo, &hi)
		} else {
			f(i, &hi, &lo)
		}
	}

	return uint64(hi-a[1])<<32 | uint64(lo-a[0])
}

/************** end encrypt name ****************/

func load(n uint64) {
	lo := uint32(n)
	hi := uint32(n >> 32)
	var x, y, t uint32

	f := func(i int, k *uint32) {
		t = b[i] + x + y
		a[i] = t<<3 | t>>29
		x = a[i]
		t = *k + x + y
		y = (x + y) & 31
		*k = t<<uint32(y) | t>>uint32(32-y)
		y = *k
	}

	for i := 0; i < 26; i++ {
		if i%2 == 0 {
			f(i, &lo)
		} else {
			f(i, &hi)
		}
	}

	g := func(i int, k *uint32) {
		t = a[i] + x + y
		a[i] = t<<3 | t>>29
		x = a[i]
		t = *k + x + y
		y = (x + y) & 31
		*k = t<<uint32(y) | t>>uint32(32-y)
		y = *k
	}

	for j := 0; j < 2; j++ {
		for i := 0; i < 26; i++ {
			if i%2 == 0 {
				g(i, &lo)
			} else {
				g(i, &hi)
			}
		}
	}
}
