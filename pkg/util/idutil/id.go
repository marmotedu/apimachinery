// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package idutil

import (
	"crypto/rand"

	"github.com/marmotedu/apimachinery/pkg/util/iputil"
	"github.com/marmotedu/apimachinery/pkg/util/stringutil"
	"github.com/sony/sonyflake"
	hashids "github.com/speps/go-hashids"
)

// Defiens alphabet.
const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

// nolint: gochecknoglobals // no need
var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		ip := iputil.GetLocalIP()

		return uint16([]byte(ip)[2])<<8 + uint16([]byte(ip)[3]), nil
	}

	sf = sonyflake.NewSonyflake(st)
}

// GetIntID returns uint64 uniq id.
func GetIntID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

// GetUUID returns id format like: B6BZVN3mOPvx.
func GetUUID(prefix string) string {
	id := GetIntID()
	hd := hashids.NewData()

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(i)
}

// GetUUID36 returns id format like: 300m50zn91nwz5.
func GetUUID36(prefix string) string {
	id := GetIntID()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(i)
}

func randString(letters string, n int) string {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letters)
	// fill output
	for pos := range output {
		// get random item
		random := randomness[pos]

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letters[randomPos]
	}

	return string(output)
}

// NewSecret returns a secretID.
func NewSecretID() string {
	return randString(Alphabet62, 36)
}

// NewSecretKey returns a secretKey or password.
func NewSecretKey() string {
	return randString(Alphabet62, 32)
}
