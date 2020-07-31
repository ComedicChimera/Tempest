package models

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type Nonce struct {
	gorm.Model
	Value     int       `gorm:"column:value"`
	Timestamp time.Time `gorm:"column:timestamp"`
	Addr      string    `gorm:"column:addr"`
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

var rnd *rand.Rand

func init() {
	var src cryptoSource
	rnd = rand.New(src)
}

func NewNonce(addr string) (int, bool) {
	nm := &Nonce{Addr: addr, Value: rnd.Int(), Timestamp: time.Now()}
	db.Create(nm)

	if nm.ID <= 0 {
		return 0, false
	}

	return nm.Value, true
}

const nonceLifetime float64 = 2.0 // 2 minutes

func GetAddrNonce(addr string) (int, bool) {
	nm := &Nonce{}

	if err := db.Table("nonces").Where("addr = ?", addr).First(nm).Error; err != nil {
		return 0, false
	}

	db.Delete(nm)
	return nm.Value, time.Since(nm.Timestamp).Minutes() <= nonceLifetime
}
