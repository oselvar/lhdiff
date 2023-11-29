package lhdiff

import (
	"fmt"

	"leb.io/hashland/spooky"
	// "github.com/tildeleb/hashland/spooky"
)

func ExampleFasterSimHash() {
	h1 := FasterSimHash("the quick brown fox")
	h2 := FasterSimHash("the quick brown fox jumped")
	h3 := FasterSimHash("and now over to something completely different")

	fmt.Println(Hamming(h1, h2))
	fmt.Println(Hamming(h1, h3))
	fmt.Println(Hamming(h2, h3))

	// Output:
	// 19
	// 21
	// 8
}

func ExampleFasterSimHash2() {
	h1 := FasterSimHash2("the quick brown fox")
	h2 := FasterSimHash2("the quick brown fox jumped")
	h3 := FasterSimHash2("and now over to something completely different")

	fmt.Println(Hamming(h1, h2))
	fmt.Println(Hamming(h1, h3))
	fmt.Println(Hamming(h2, h3))

	// Output:
	// 13
	// 27
	// 24
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
