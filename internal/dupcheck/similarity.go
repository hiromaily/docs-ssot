package dupcheck

import (
	"math"
)

type vector map[string]float64

// buildTFIDF computes a TF-IDF weighted, unit-normalised vector for each chunk.
func buildTFIDF(chunks []Chunk) []vector {
	tfList := make([]map[string]float64, len(chunks))
	df := map[string]float64{}

	for i, c := range chunks {
		tf := map[string]float64{}
		seen := map[string]bool{}

		for _, tok := range c.Tokens {
			tf[tok]++
			if !seen[tok] {
				df[tok]++
				seen[tok] = true
			}
		}
		tfList[i] = tf
	}

	n := float64(len(chunks))
	idf := make(map[string]float64, len(df))
	for term, freq := range df {
		idf[term] = math.Log((n+1.0)/(freq+1.0)) + 1.0
	}

	out := make([]vector, len(chunks))
	for i, tf := range tfList {
		vec := make(vector, len(tf))
		var norm float64

		for term, freq := range tf {
			weight := freq * idf[term]
			vec[term] = weight
			norm += weight * weight
		}

		if norm > 0 {
			sq := math.Sqrt(norm)
			for term, weight := range vec {
				vec[term] = weight / sq
			}
		}
		out[i] = vec
	}

	return out
}

// cosine computes the cosine similarity between two unit-normalised vectors.
func cosine(a, b vector) float64 {
	if len(a) > len(b) {
		a, b = b, a
	}

	var dot float64
	for term, av := range a {
		if bv, ok := b[term]; ok {
			dot += av * bv
		}
	}
	return dot
}
