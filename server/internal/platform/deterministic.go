package platform

// IDs is a deterministic ID source for tests and fixtures.
type IDs struct {
	values []string
	index  int
}

func NewIDs(values ...string) *IDs { return &IDs{values: append([]string(nil), values...)} }
func (g *IDs) Next() (string, error) {
	if g.index >= len(g.values) {
		return NewID()
	}
	value := g.values[g.index]
	g.index++
	return value, nil
}
