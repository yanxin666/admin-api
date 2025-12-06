package builder

import (
	"context"
)

type WhiteListChecker struct {
	whiteList []string
}

func (c *WhiteListChecker) SetRule(ctx context.Context, input interface{}) bool {
	whiteList, ok := input.([]string)
	if !ok {
		return false
	}
	c.whiteList = whiteList
	return true
}

func (c *WhiteListChecker) CheckRuleValid(param *Params) bool {
	if len(c.whiteList) == 0 {
		return false
	}

	return true
}

func (c *WhiteListChecker) Do(ctx context.Context, params *Params) (bool, int) {
	return true, 1
}
