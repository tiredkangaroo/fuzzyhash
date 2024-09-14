package fuzzyhash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"math"
	"slices"
	"strconv"
	"strings"
)

// FuzzyHash represents a hash generated fuzzily.
type FuzzyHash struct {
	// C is populated as the output of the fuzzy
	// hash function.
	C []byte
	// T is the number of bytes generated from the
	// input bytes. The higher this number is the less
	// fuzzy the hash is.
	T int
}

// roughMedian calculates the rough median of an array of numbers.
func roughMedian[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](s []T) T {
	var numbers []float64
	for _, num := range s {
		numbers = append(numbers, float64(num))
	}
	slices.Sort(numbers)
	return T(numbers[int(math.Floor(float64(len(numbers))/2.0))])
}

// ExtractBytes takes in input bytes (k) and an integer (s) to extract
// an s number of bytes from the input bytes.
func ExtractBytes(k []byte, s int) ([]byte, error) {
	if s == len(k) {
		return k, nil
	}
	if s > len(k) {
		return []byte{}, fmt.Errorf("s may not be larger than length of k")
	}
	kh := []uint8(hex.EncodeToString(k))              // encode []byte as hexadecimal
	t := len(kh) / s                                  // get the step (how many bytes to one byte)
	kha := roughMedian(kh)                            // calculate the rough median
	padding := slices.Repeat([]uint8{kha}, len(kh)%t) // pad it so the step (padding is the rough median) can be used in a loop properly
	key := append(kh, padding...)                     // add the padding
	var uc []uint8                                    // uc are the extracted bytes
	for i := range len(key) / t {
		buffer := key[i*t : i*t+t]
		s := roughMedian(buffer) // get the median of the current step
		uc = append(uc, s)       // add the median as an extracted byte
	}
	return uc, nil // return the extracted bytes
}

// HashWith takes in input bytes (k) and integer (s) to measure the fuzzyness.
// The larger s is, the more fuzzy the hash will be. The hashing system used
// will be h.
func HashWith(h hash.Hash, k []byte, s int) (*FuzzyHash, error) {
	output, err := ExtractBytes(k, s)
	if err != nil {
		return nil, err
	}
	fuzzyhash := new(FuzzyHash)
	fuzzyhash.T = s
	_, err = h.Write(output)
	if err != nil {
		return nil, err
	}
	fuzzyhash.C = h.Sum(nil)
	return fuzzyhash, nil
}

func MustHash(f *FuzzyHash, e error) *FuzzyHash {
	if e != nil {
		panic(e)
	}
	return f
}

// HashSHA1 is equal to call to HashWith with the first parameter
// being sha1.New()
func HashSHA1(k []byte, s int) (*FuzzyHash, error) {
	return HashWith(sha1.New(), k, s)
}

// Hash256 is equal to call to HashWith with the first parameter
// being sha256.New()
func Hash256(k []byte, s int) (*FuzzyHash, error) {
	return HashWith(sha256.New(), k, s)
}

// Hash512 is equal to call to HashWith with the first parameter
// being sha512.New()
func Hash512(k []byte, s int) (*FuzzyHash, error) {
	return HashWith(sha512.New(), k, s)
}

// String returns the output bytes in hexadecimal form.
func (f FuzzyHash) String() string {
	return hex.EncodeToString(f.C)
}

// MarshalString marshals the FuzzyHash into a string.
func (f FuzzyHash) MarshalString() string {
	c := hex.EncodeToString(f.C)
	return fmt.Sprintf("%s:%d", c, f.T)
}

// UnmarshalString unmarshals a string (generated by MarshalString)
// into a FuzzyHash. It may return an error if the string is invalid.
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

// Equal checks for the equality of the output bytes
// of two fuzzyhashes.
func Equal(f1 *FuzzyHash, f2 *FuzzyHash) bool {
	return f1.String() == f2.String()
}
