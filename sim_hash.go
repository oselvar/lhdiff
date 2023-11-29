package lhdiff

import (
	"math/bits"

	"leb.io/hashland/siphash"
	"leb.io/hashland/spooky"
)

// https://ferd.ca/simhashing-hopefully-made-simple.html
func FasterSimHash(s string) uint64 {
	bitLength := 64
	shingleSize := 2
	seed := uint64(0xfedcba)

	bits := make([]uint64, bitLength)
	shingles := GenerateShingles(s, shingleSize)
	for _, shingle := range shingles {
		_, hash := spooky.SpookyHash128(shingle, seed, seed)

		for i := bitLength; i >= 1; i-- {
			if ((hash >> (bitLength - i)) & 1) == 1 {
				bits[i-1]++
			} else {
				bits[i-1]--
			}
		}
	}

	var simHash uint64

	one := uint64(1)
	for i := bitLength; i >= 1; i-- {
		if bits[i-1] > 0 {
			simHash |= one
		}
		one = one << 1
	}
	return simHash
}

func FasterSimHash2(s string) uint64 {
	shingleSize := 2
	// seed := uint64(0xfedcba)

	// bits := make([]uint64, bitLength)
	var signs [64]int64
	shingles := GenerateShingles(s, shingleSize)
	for _, shingle := range shingles {
		// _, h := spooky.SpookyHash128(shingle, seed, seed)
		h := siphash.Hash(0, 0, shingle)

		for i := 0; i < 64; i++ {
			negate := int(h) & 1
			// if negate is 1, we will negate '-1', below
			r := (-1 ^ -negate) + negate
			signs[i] += int64(r)
			h >>= 1
		}
	}

	var shash uint64

	// TODO: can probably be done with SSE?
	for i := 63; i >= 0; i-- {
		shash <<= 1
		shash |= uint64(signs[i]>>63) & 1
	}

	return shash
}

func GenerateShingles(input string, tokenLength int) [][]byte {
	// Create a slice to store the shingles
	// var shingles []string
	var tokens [][]byte

	var token []byte
	// Iterate over each character in the input string
	for i := 0; i < len(input); i++ {
		if token == nil {
			token = make([]byte, 0, tokenLength)
		}
		// Skip whitespace characters
		if input[i] != ' ' {
			token = append(token, input[i])
			if len(token) == tokenLength {
				tokens = append(tokens, token)
				token = nil
			}
		}
	}
	if token != nil {
		for len(token) < tokenLength {
			token = append(token, '_')
		}
		tokens = append(tokens, token)
	}

	return tokens
}

func Hamming(b1, b2 uint64) int {
	return bits.OnesCount64(b1 ^ b2)
}

func Hamming2(v1, v2 uint64) int {
	v := v1 ^ v2
	var c int
	for c = 0; v != 0; c++ {
		v &= v - 1
	}

	return c
}
