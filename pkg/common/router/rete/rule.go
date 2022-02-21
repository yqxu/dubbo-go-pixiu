package rete

type Rule interface {
}

type RuleImpl struct {
	Rule
	id int64
}

type SimpleRule struct {
	RuleImpl
	operator    Operator
	staticValue string
}

const (
	and = "and"
	or  = "or"
)

type CombineRule struct {
	RuleImpl
	logical    string
	childRules []Rule
}
