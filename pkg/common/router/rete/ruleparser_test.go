package rete

import (
	"go/parser"
	"testing"
)

func TestParse(t *testing.T) {
	parser.ParseExpr("ruleStr.str.a[2].c[4] >= 100")

}
