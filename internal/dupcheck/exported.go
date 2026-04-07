package dupcheck

// Vector is a TF-IDF weighted, unit-normalised term vector.
type Vector = vector

// BuildTFIDF computes TF-IDF weighted, unit-normalised vectors for the given chunks.
// Chunks must have their Tokens field populated (use Tokenize).
func BuildTFIDF(chunks []Chunk) []Vector {
	return buildTFIDF(chunks)
}

// Cosine computes the cosine similarity between two unit-normalised vectors.
// Returns a value in [0.0, 1.0] where 1.0 means identical.
func Cosine(a, b Vector) float64 {
	return cosine(a, b)
}

// Tokenize splits text into tokens suitable for TF-IDF comparison.
func Tokenize(s string) []string {
	return tokenize(s)
}
