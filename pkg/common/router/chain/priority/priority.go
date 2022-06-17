package priority

const (
	URLMatch    = "path"
	HeaderMatch = "headers"
)

var priorityMap = map[string]int{
	URLMatch:    0,
	HeaderMatch: 1,
}

func priorityOf(matcherCode string) int {
	return priorityMap[matcherCode]
}

func hasLowerPriorityConf(current string, configKeys []string) bool {
	//TODO
	return true
}
