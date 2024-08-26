package fuzzyhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type FuzzyHash struct {
	C []byte
	T int
}

func roughMedian[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](s []T) T {
	var numbers []float64
	for _, num := range s {
		numbers = append(numbers, float64(num))
	}
	slices.Sort(numbers)
	return T(numbers[int(math.Floor(float64(len(numbers))/2.0))])
}

func Hash(k []byte, t float64) *FuzzyHash {
	tolerance := int(math.Ceil((12 / float64(len(k)) * float64(len(k)) * t)))
	kh := []uint8(hex.EncodeToString(k))
	kha := roughMedian(kh)
	padding := slices.Repeat([]uint8{kha}, len(kh)%tolerance)
	key := append(kh, padding...)
	fuzzyHash := &FuzzyHash{T: tolerance}
	var uc []uint8
	for i := range len(key) / tolerance {
		buffer := key[i*tolerance : i*tolerance+tolerance]
		s := roughMedian(buffer)
		uc = append(uc, s)
	}
	hasher := sha256.New()
	hasher.Write(uc)
	fuzzyHash.C = hasher.Sum(nil)
	return fuzzyHash
}

func (f FuzzyHash) String() string {
	return hex.EncodeToString(f.C)
}

func (f FuzzyHash) MarshalString() string {
	c := hex.EncodeToString(f.C)
	return fmt.Sprintf("%s:%d", c, f.T)
}

func UnmarshalString(fh string) (*FuzzyHash, error) {
	nf := new(FuzzyHash)
	parts := strings.Split(fh, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("string cannot be made into fuzzyhash: its formatting is invalid.")
	}
	ce := parts[0]
	ts := parts[1]

	c, err := hex.DecodeString(ce)
	if err != nil {
		return nil, fmt.Errorf("string cannot be made into fuzzyhash: specified c value is not hex encoded")
	}
	nf.C = c

	t, err := strconv.Atoi(ts)
	if err != nil {
		return nil, fmt.Errorf("string cannot be made into fuzzyhash: specified t value is not an integer")
	}
	nf.T = t
	return nf, nil
}

func Equal(f1 *FuzzyHash, f2 *FuzzyHash) bool {
	return f1.String() == f2.String()
}
