package rete

import (
	_ "go/ast"
	_ "go/token"
)

type Path struct {
}

func parse(ruleStr string) Rule {
	//TODO: 确认配置文件使用的语法 决定ast相关技术栈选型 暂定go/ast
	if isSimpleRule(ruleStr) {
		return SimpleRule{operator: GetOperator("=")}
	}
	return CombineRule{logical: and}
}

func isSimpleRule(ruleStr string) bool {
	return true
}
