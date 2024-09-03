package printer

var Discard = discard{}

type discard struct{}

func (d discard) Print(message string) {}
