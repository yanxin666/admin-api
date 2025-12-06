package builder

import "context"

type CheckerName string

//type CheckerType string

type CheckerUnit struct {
	//CheckerType
	CheckerName
	Checker
}

type GetCheckerFunc func() Checker

type Checker interface {
	// SetRule 设置checker规则
	SetRule(ctx context.Context, filterConf interface{}) bool
	// CheckRuleValid 检查当前filter是否有效
	CheckRuleValid(param *Params) bool
	// Do 执行filter
	Do(ctx context.Context, param *Params) (bool, int)
}

const (
	CheckerWhiteList      = "white_list"
	CheckerOpenStatus     = "open_status"
	CheckerBigDataTag     = "big_data_tag"
	CheckerFenceInfo      = "fence_info"
	CheckerDatePeriod     = "date_period"
	CheckerWeekList       = "week_list"
	CheckerTimePeriodList = "time_period_list"
)

type Params struct {
	AccessKeyID   int32
	AppVersion    string
	Pid           int64
	Uid           int64
	Phone         string
	FromLat       float64
	FromLng       float64
	ToLat         float64
	ToLng         float64
	OrderType     int32
	DepartureTime int64
	MapType       string
	BigDataTag    []string
}

var CheckerMap = map[string]GetCheckerFunc{
	CheckerWhiteList: func() Checker { return &WhiteListChecker{} },
	//CheckerOpenStatus: func() Checker { return &OpenStatusChecker{} },
	//CheckerBigDataTag: func() Checker { return &TagChecker{} },
	//CheckerFenceInfo: func() Checker { return &FenceChecker{} },
	//CheckerDatePeriod: func() Checker { return &DateChecker{} },
	//CheckerWeekList: func() Checker { return &WeekdayChecker{} },
	//CheckerTimePeriodList: func() Checker { return &TimeChecker{} },
}

func LoadChecker(ctx context.Context, checkerName string, checkerConf interface{}, param *Params) *CheckerUnit {
	var getCheckerFunc GetCheckerFunc
	var checker Checker
	if getFunc, ok := CheckerMap[checkerName]; ok {
		getCheckerFunc = getFunc
	}

	if getCheckerFunc == nil {
		return nil
	}

	checker = getCheckerFunc()
	if checker == nil {
		return nil
	}

	if checker.SetRule(ctx, checkerConf) && checker.CheckRuleValid(param) {
		cUnit := &CheckerUnit{}
		cUnit.CheckerName = CheckerName(checkerName)
		cUnit.Checker = checker
		return cUnit
	}

	return nil
}

func ExecChecker(ctx context.Context, checkerList []*CheckerUnit, params *Params) (bool, int) {
	if len(checkerList) == 0 {
		return false, 0
	}

	for _, unit := range checkerList {
		if unit != nil && unit.Checker != nil {
			if isHit, reasonTag := unit.Checker.Do(ctx, params); !isHit {
				return false, reasonTag
			} else if unit.CheckerName == CheckerWhiteList {
				return true, 1
			}
		}
	}

	return true, 1
}
