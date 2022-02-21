package rete

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

	urlMatch = "urlMatch"
)

type Operator struct {
	stmt     string
	priority int
}

var Operators = map[string]Operator{
	urlMatch: {stmt: urlMatch, priority: 10000},
	equals:   {stmt: equals, priority: 9000},
	in:       {stmt: in, priority: 1001},
	contains: {stmt: contains, priority: 1000},

	bigger:         {stmt: bigger, priority: 1},
	biggerOrEquals: {stmt: biggerOrEquals, priority: 2},
	less:           {stmt: less, priority: 3},
	lessOrEquals:   {stmt: lessOrEquals, priority: 4},

	notEquals:   {stmt: notEquals, priority: -1},
	notContains: {stmt: notContains, priority: -2},
	notIn:       {stmt: notIn, priority: -3},
}

//priority : unique batch match result have priority

func GetOperator(name string) Operator {
	return Operators[name]
}
