package alphaRete

const (
	equals    = "=="
	notEquals = "!="

	in    = "in"
	notIn = "notIn"

	contains    = "contains"
	notContains = "notContains"

	bigger         = ">"
	biggerOrEquals = ">="
	less           = "<"
	lessOrEquals   = "<="
)

type operator struct {
	stmt     string
	priority int
}
