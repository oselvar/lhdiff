package lhdiff

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	levenshtein "github.com/ka-weihe/fast-levenshtein"
)

type LineInfo struct {
	lineNumber     int
	content        string
	context        string
	contentSimHash uint64
	contextSimHash uint64
}

type LinePair struct {
	left  *LineInfo
	right *LineInfo
}

func MakeLineInfos(lineNumbers []int, lines []string, contextSize int) []*LineInfo {
	lineInfos := make([]*LineInfo, len(lineNumbers))
	for i, lineNumber := range lineNumbers {
		lineInfos[i] = MakeLineInfo(lineNumber, lines, contextSize)
	}
	return lineInfos
}

func MakeLineInfo(lineNumber int, lines []string, contextSize int) *LineInfo {
	content := lines[lineNumber]
	context := getContext(lineNumber, lines, contextSize)

	contentSimHash := FasterSimHash2(content)
	contextSimHash := FasterSimHash2(context)

	lineInfo := &LineInfo{
		lineNumber:     lineNumber,
		context:        context,
		content:        content,
		contentSimHash: contentSimHash,
		contextSimHash: contextSimHash,
	}
	return lineInfo
}

func (linePair *LinePair) contentHamming() int {
	return Hamming2(linePair.left.contentSimHash, linePair.right.contentSimHash)
}

func (linePair *LinePair) contextHamming() int {
	return Hamming2(linePair.left.contextSimHash, linePair.right.contextSimHash)
}

func (linePair *LinePair) contentNormalizedLevenshteinSimilarity() float64 {
	distance := levenshtein.Distance(linePair.left.content, linePair.right.content)
	normalizedLevenhsteinDistance := float64(distance) / math.Max(float64(len(linePair.left.content)), float64(len(linePair.right.content)))
	return 1 - normalizedLevenhsteinDistance
}

func (linePair *LinePair) contextTfIdfCosineSimilarity() float64 {
	return TfIdfCosineSimilarity(linePair.left.context, linePair.right.context)
}

func (linePair *LinePair) combinedSimilarity() float64 {
	contentSimilarity := linePair.contentNormalizedLevenshteinSimilarity()
	if contentSimilarity <= 0.5 {
		return 0.0
	}
	contextSimilarity := linePair.contextTfIdfCosineSimilarity()
	return ContentSimilarityFactor*contentSimilarity + ContextSimilarityFactor*contextSimilarity
}

func (linePair *LinePair) combinedHammingSimilarity() float32 {
	contentHamming := linePair.contentHamming()
	contextHamming := linePair.contextHamming()
	combined := ContentSimilarityFactor*(float32(contentHamming)/32) + ContextSimilarityFactor*(float32(contextHamming)/32)
	fmt.Printf("combined: %f\n", combined)
	return combined
}

type ByCombinedSimilarity []LinePair

func (a ByCombinedSimilarity) Len() int { return len(a) }
func (a ByCombinedSimilarity) Less(i, j int) bool {
	return a[j].combinedSimilarity() < a[i].combinedSimilarity()
}
func (a ByCombinedSimilarity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

var /* const */ brackets = regexp.MustCompile("^[{()}]$")

// getContext returns a string consisting of (up to) contextSize context lines above and below lineIndex.
// a line is considered to be a context line if it is not an "insignificant" line, i.e. either blank
// or just a curly brace or parenthesis (whitespace trimmed).
func getContext(lineNumber int, lines []string, contextSize int) string {
	var context []string

	i := lineNumber - 1

	for j := 0; i >= 0 && j < contextSize; {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if len(trimmed) != 0 && !brackets.MatchString(trimmed) {
			context = append([]string{line}, context...)
			j++
		}
		i--
	}

	i = lineNumber + 1
	for j := 0; i < len(lines) && j < contextSize; {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if len(trimmed) != 0 && !brackets.MatchString(trimmed) {
			context = append(context, line)
			j++
		}
		i++
	}

	return strings.Join(context, "")
}
