package lhdiff

import (
	"sort"

	"github.com/sourcegraph/go-diff/diff"
)

func Algo1(
	fileDiff *diff.FileDiff,
	leftLines []string,
	rightLines []string,
	contextSize int,
	mappedByRightLineNumber *map[int]bool,
	linePairByLeftLineNumber *map[int]LinePair) {
	unchangedDiffPairs, leftLineNumbers, rightLineNumbers := LineNumbersFromDiff(fileDiff, leftLines, rightLines, contextSize)
	for _, unchangedDiffPair := range unchangedDiffPairs {
		(*linePairByLeftLineNumber)[unchangedDiffPair.left.lineNumber] = unchangedDiffPair
		(*mappedByRightLineNumber)[unchangedDiffPair.right.lineNumber] = true
	}

	leftLineInfos := MakeLineInfos(leftLineNumbers, leftLines, contextSize)
	rightLineInfos := MakeLineInfos(rightLineNumbers, rightLines, contextSize)

	// TODO: We have combinatorial explosion here....
	// See section D in the paper about simhash
	// We need to compute that here.
	// See HDiffSHMatching.match - line 287-293
	// Maybe do this in parallel?
	for _, rightLineInfo := range rightLineInfos {
		var similarPairCandidates []LinePair
		for _, leftLineInfo := range leftLineInfos {
			// float k = CONTEXT_SIMILARITY_FACTOR * ((float) i / 32.0F) + CONTENT_SIMILARITY_FACTOR * ((float) m / 32.0F);
			pair := LinePair{
				left:  leftLineInfo,
				right: rightLineInfo,
			}
			// hammingSimilarity := pair.combinedHammingSimilarity()
			// if hammingSimilarity > SimilarityThreshold {
			// fmt.Println("hi:", hammingSimilarity)
			similarPairCandidates = append(similarPairCandidates, pair)
			// } else {
			// fmt.Println("lo:", hammingSimilarity)
			// }
		}
		sort.Sort(ByCombinedSimilarity(similarPairCandidates))
		if len(similarPairCandidates) > 0 {
			mostSimilarPair := similarPairCandidates[0]
			if mostSimilarPair.combinedSimilarity() > SimilarityThreshold {
				(*linePairByLeftLineNumber)[mostSimilarPair.left.lineNumber] = mostSimilarPair
				(*mappedByRightLineNumber)[mostSimilarPair.right.lineNumber] = true
			}
		}
	}
}
