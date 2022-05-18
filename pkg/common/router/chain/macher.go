package chain

import stdHttp "net/http"

type MatchPathNode struct {
	MatcherName string
	Path        string
}

type MatchPath struct {
	Path []MatchPathNode
}

type Chain interface {
	GetRootMatcherNode() MatchNode
	AddAtTop(node MatchNode)
	AppendAtEnd(node MatchNode)
	Top() MatchNode
	End() MatchNode
}

type Matcher interface {
	DoMatch(r *stdHttp.Request, path MatchPath) MatchResult
	getCode() string
}

type RouteConfigSet interface {
	FindRouteConfig(ruleKey string) interface{}
	RouteConfigPut(ruleKey string, rule interface{})
	RouteConfigRemove(ruleKey string)
}

func EndOfChainSuccessResult(path MatchPath, bizData interface{}) MatchResult {
	return MatchResult{EndOfChain: true, MatchSuccess: true, BizData: bizData, MatchPath: path}
}

func MidOfChainSuccessResult(path MatchPath, bizData interface{}) MatchResult {
	return MatchResult{EndOfChain: false, MatchSuccess: true, BizData: bizData, MatchPath: path}
}

func EndOfChainFailedResult() MatchResult {
	return MatchResult{EndOfChain: true, MatchSuccess: false}
}

func MidOfChainFailedResult() MatchResult {
	return MatchResult{EndOfChain: false, MatchSuccess: false}
}

type MatchResult struct {
	EndOfChain   bool
	MatchSuccess bool
	MatchPath    MatchPath
	BizData      interface{}
}

type MatchNode interface {
	Matcher
	//Next get the next node after this.
	Next() MatchNode
	//Append to append a MatchNode after this.
	Append(node MatchNode) MatchNode
}

type ParamMapping interface {
	GetUseful(r *stdHttp.Request) interface{}
}

//AlwaysTrue The match Node append at the end of chain.
type AlwaysTrue struct {
	MatchNode
}

func BuildAlwaysTrue() AlwaysTrue {
	return AlwaysTrue{}
}

func (matcher *AlwaysTrue) DoMatch(r *stdHttp.Request, path MatchPath) bool {
	return true
}

func (matcher *AlwaysTrue) next() MatchNode {
	return nil
}

func (matcher *AlwaysTrue) Append(node MatchNode) {
	return
}

func (matcher *AlwaysTrue) Config() {
	return
}
func (matcher *AlwaysTrue) PutConfig(conf interface{}) {
	return
}
