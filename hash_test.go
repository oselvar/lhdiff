package lhdiff

import (
	"fmt"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/token/shingle"
	"leb.io/hashland/spooky"
	// "github.com/tildeleb/hashland/spooky"
)

func FasterSimHash(s string) uint64 {
	hashSize := 64
	shingleSize := 2
	seed := uint64(0xfedcba)

	vector := make([]uint64, hashSize)
	shingles := GenerateShingles(s, shingleSize)
	for _, shingle := range shingles {
		_, shingleHashLow := spooky.SpookyHash128(shingle, seed, seed)
		// shingleHash :=
		for c := 0; c < hashSize; c++ {
			if shingleHashLow&(1<<c) == 0 {
				vector[c]--
			} else {
				vector[c]++
			}
		}
	}

	var result uint64

	for s := 0; s < hashSize; s++ {
		if vector[s] > 0 {
			result |= 1 << s
		} else {
			result |= 0 << s
		}
	}
	return result
}

func ExampleFasterSimHash() {
	hash := FasterSimHash("the quick brown fox")
	fmt.Println(hash)

	// Output:
	// 17795968136258301881
}

func ExampleSpookyHashShort() {
	seed := uint64(0xfedcba)
	a, b := spooky.SpookyHash128([]byte("hello"), seed, seed)
	fmt.Println(a)
	fmt.Println(b)

	// Output:
	// 1916231864266175103
	// 13411944305656060552
}

func ExampleNewShingleFilter() {
	input := analysis.TokenStream{
		&analysis.Token{
			Term: []byte("the"),
		},
		&analysis.Token{
			Term: []byte("quick"),
		},
		&analysis.Token{
			Term: []byte("brown"),
		},
		&analysis.Token{
			Term: []byte("fox"),
		},
	}
	shingleFilter := shingle.NewShingleFilter(2, 2, false, " ", "_")
	actual := shingleFilter.Filter(input)
	for _, token := range actual {
		fmt.Println(string(token.Term))
	}

	// Output:
	// the quick
	// quick brown
	// brown fox
}

func ExampleGenerateShingles_without_padding() {
	shingles := GenerateShingles("the quick brown fox", 2)
	for _, shingle := range shingles {
		fmt.Println(string(shingle))
	}

	// Output:
	// th
	// eq
	// ui
	// ck
	// br
	// ow
	// nf
	// ox
}

func ExampleGenerateShingles_with_padding() {
	shingles := GenerateShingles("the quick brown fox", 5)
	for _, shingle := range shingles {
		fmt.Println(string(shingle))
	}

	// Output:
	// thequ
	// ickbr
	// ownfo
	// x____
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
