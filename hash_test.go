package lhdiff

import (
	"fmt"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/token/shingle"
	// "github.com/tildeleb/hashland/spooky"
)

// func ExampleNew_spooky() {
// 	seed := uint64(0xfedcba)
// 	a, b := spooky.SpookyHashShort([]byte("hello"), seed, seed)
// 	fmt.Println(a)
// 	fmt.Println(b)
// 	// fmt.Println(hash.Sum([]byte("hello")))

// 	// Output:
// 	// 1916231864266175103
// 	// 13411944305656060552
// }

func ExampleNew_shingle() {
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
