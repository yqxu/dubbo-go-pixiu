package trie

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/chain"
	stdHttp "net/http"
)

type UrlMatcher struct {
	chain.MatchNode
	chain.RouteConfigSet

	BelongTo chain.Chain
	NextNode chain.MatchNode
	Conf     Trie
}

type UrlParamMapping struct {
	chain.ParamMapping
}

func getUseful(r *stdHttp.Request) string {
	return r.URL.Host
}

func (matcher UrlMatcher) Append(matchNode chain.MatchNode) {
	matcher.NextNode = matchNode
}

func (matcher UrlMatcher) Next() chain.MatchNode {
	return matcher.NextNode
}

func (matcher UrlMatcher) DoMatch(r *stdHttp.Request, path chain.MatchPath) {
	//Always the first node of chain so ,get config by MatchPath is unnecessary.
	matcher.Conf.MatchForChain(r)
}

func BuildUrlMatchPath(path string) chain.MatchPath {
	return chain.MatchPath{Path: []chain.MatchPathNode{chain.MatchPathNode{
		MatcherName: "URLMatch",
		Path:        path,
	}}}
}
