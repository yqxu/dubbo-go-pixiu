package sequential

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/chain"
	"net/http"
)

type HeaderMatcher struct {
}

func (h HeaderMatcher) FindRouteConfig(rule interface{}) interface{} {
	panic("implement me")
}

func (h HeaderMatcher) RouteConfigPut(rule interface{}) {
	panic("implement me")
}

func (h HeaderMatcher) RouteConfigRemove(ruleKey string) {
	panic("implement me")
}

func (h HeaderMatcher) DoMatch(r *http.Request, path chain.MatchPath) chain.MatchResult {
	panic("implement me")
}

func (h HeaderMatcher) getCode() string {
	panic("implement me")
}

func (h HeaderMatcher) Next() chain.MatchNode {
	panic("implement me")
}

func (h HeaderMatcher) Append(node chain.MatchNode) chain.MatchNode {
	panic("implement me")
}

func BuildHeaderMatcher() chain.RouteConfigSet {
	return HeaderMatcher{}
}
