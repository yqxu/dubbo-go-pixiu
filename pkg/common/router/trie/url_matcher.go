package trie

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/chain"
	stdHttp "net/http"
)

type UrlMatcher struct {
	chain.RouteConfigSet
	chain.MatchNode

	BelongTo chain.Chain
	NextNode chain.MatchNode
	Conf     Trie
}

func (matcher UrlMatcher) getCode() string {
	return "url"
}

func (matcher UrlMatcher) RouteConfigRemove(ruleKey string) {
	matcher.Conf.Remove(ruleKey)
}

type UrlParamMapping struct {
	chain.ParamMapping
}

func getUseful(r *stdHttp.Request) string {
	return r.URL.Host
}

func (matcher UrlMatcher) Append(matchNode chain.MatchNode)  chain.MatchNode{
	matcher.NextNode = matchNode
	return matchNode
}

func (matcher UrlMatcher) Next() chain.MatchNode {
	return matcher.NextNode
}

func (matcher UrlMatcher) DoMatch(r *stdHttp.Request, path chain.MatchPath) chain.MatchResult{
	//Always the first node of chain so ,get config by MatchPath is unnecessary.
	return matcher.Conf.MatchForChain(r)
}

func BuildUrlMatchPath(path string) chain.MatchPath {
	return chain.MatchPath{Path: []chain.MatchPathNode{chain.MatchPathNode{
		MatcherName: "URLMatch",
		Path:        path,
	}}}
}

func (matcher UrlMatcher) FindRouteConfig(rule interface{}) interface{}{
	return _,_,_,_ = matcher.Conf.Get(rule)
}
func (matcher UrlMatcher) RouteConfigPut(rule interface{}){
	matcher.Conf.Put()
}

func RouteConfigRemove(ruleKey string){

}




