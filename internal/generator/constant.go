package generator

// ConstantGenerator is a Generator implementation that always returns the same constant value.
// It is not suitable for production use where unique identifiers are required,
// but is valuable in tests and examples where deterministic output is needed.
type ConstantGenerator struct {
	value string // value is the constant string returned by Generate.
}

func NewConstantGenerator(value string) *ConstantGenerator {
	return &ConstantGenerator{
		value: value,
	}
}

func (g *ConstantGenerator) Generate() (string, error) {
	return g.value, nil
}
