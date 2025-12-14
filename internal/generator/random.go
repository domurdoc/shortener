package generator

import "github.com/domurdoc/shortener/internal/utils"

// RandomGenerator is a Generator implementation that produces random short IDs.
// It generates strings of a fixed length by selecting characters from a specified charset.
// Each call to Generate returns a new random string, suitable for use as a short URL key.
type RandomGenerator struct {
	charset string
	length  int
}

func NewRandomGenerator(charset string, length int) *RandomGenerator {
	return &RandomGenerator{charset: charset, length: length}
}

func (g *RandomGenerator) Generate() (string, error) {
	return utils.GenerateRandomString(g.charset, g.length)
}
