package lcnc

type OriginContext struct {
	Origin   Origin
	ForBlock bool
	Query    bool
	Input    bool
	Output   bool
}

type Origin string

const (
	OriginInvalid  Origin = "invalid"
	OriginFor      Origin = "for"
	OriginVariable Origin = "variable"
	Originfunction Origin = "function"
)
