package fuzzyhash

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"slices"
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

func Equal(f1 *FuzzyHash, f2 *FuzzyHash) bool {
	return f1.String() == f2.String()
}
