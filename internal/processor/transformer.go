package processor

// Transformer transforms the text content of a Markdown document.
// Implementations must be safe to call with an empty string and must return
// the content unchanged when no transformation applies.
type Transformer interface {
	Transform(content string) (string, error)
}

// Apply runs content through each transformer in order, returning the final result.
// The first error encountered is returned immediately.
func Apply(content string, transformers ...Transformer) (string, error) {
	for _, t := range transformers {
		result, err := t.Transform(content)
		if err != nil {
			return "", err
		}
		content = result
	}
	return content, nil
}
