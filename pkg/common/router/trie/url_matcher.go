package trie

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/chain"
	priority "github.com/apache/dubbo-go-pixiu/pkg/common/router/chain/priority"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	stdHttp "net/http"
	"strings"
	"sync"
)

type UrlMatcher struct {
	chain.RouteConfigSet
	chain.MatchNode

	BelongTo chain.Chain
	NextNode chain.MatchNode
	Conf     Trie
	confRw   sync.RWMutex
}

func (matcher UrlMatcher) getCode() string {
	return priority.URLMatch
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

func (matcher UrlMatcher) Append(matchNode chain.MatchNode) chain.MatchNode {
	matcher.NextNode = matchNode
	return matchNode
}

func (matcher UrlMatcher) Next() chain.MatchNode {
	return matcher.NextNode
}

func (matcher UrlMatcher) DoMatch(r *stdHttp.Request, path chain.MatchPath) chain.MatchResult {
	//Always the first node of chain so ,get config by MatchPath is unnecessary.
	return matcher.Conf.MatchForChain(r)
}

func BuildUrlMatchPath(path string) chain.MatchPath {
	return chain.MatchPath{Path: []chain.MatchPathNode{chain.MatchPathNode{
		MatcherName: "URLMatch",
		Path:        path,
	}}}
}

func (matcher UrlMatcher) FindRouteConfig(ruleKey string) interface{} {
	node, _, _, _ := matcher.Conf.Get(ruleKey)
	return node.GetBizInfo()
}

func (matcher UrlMatcher) RouteConfigPut(rule model.RouterMatch) {
	matcher.confRw.Lock()
	defer matcher.confRw.Unlock()
	if rule.Methods == nil {
		rule.Methods = []string{constant.Get, constant.Put, constant.Delete, constant.Post}
	}
	isPrefix := rule.Prefix != ""

	if rule.Path != "" || rule.Prefix != "" {
		for _, method := range rule.Methods {
			var key string
			if isPrefix {
				key = getTrieKey(method, method, isPrefix)
			} else {
				key = getTrieKey(method, method, isPrefix)
			}
			matcher.Conf.PutOrUpdate(key)
		}
	}

}

func RouteConfigRemove(ruleKey string) {

}

func getTrieKey(method string, path string, isPrefix bool) string {
	if isPrefix {
		if !strings.HasSuffix(path, constant.PathSlash) {
			path = path + constant.PathSlash
		}
		path = path + "**"
	}
	return stringutil.GetTrieKey(method, path)
}
