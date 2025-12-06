package define

import (
	"encoding/json"
	"muse-admin/pkg/errs"
)

type Grade struct {
	ID     int64  `json:"id"`     // id
	Name   string `json:"name"`   // 年级描述
	IsOpen bool   `json:"isOpen"` // 是否开放
}

var (
	// GradeJunior 低年级
	GradeJunior = struct {
		Preschool int64 // 学龄前
		First     int64 // 一年级
		Second    int64 // 二年级
		Third     int64 // 三年级
		Fourth    int64 // 四年级
		Fifth     int64 // 五年级
		Sixth     int64 // 六年级
	}{
		Preschool: 10,
		First:     11,
		Second:    12,
		Third:     13,
		Fourth:    14,
		Fifth:     15,
		Sixth:     16,
	}

	// GradeMiddle 中年级
	GradeMiddle = struct {
		Transition int64 // 中小衔接
		First      int64 // 初一
		Second     int64 // 初二
		Third      int64 // 中考
	}{
		Transition: 20,
		First:      21,
		Second:     22,
		Third:      23,
	}

	// GradeSenior 高年级
	GradeSenior = struct {
		First  int64 // 高一
		Second int64 // 高二
		Third  int64 // 高考
	}{
		First:  31,
		Second: 32,
		Third:  33,
	}

	// GradeMajor 专业主修年级
	GradeMajor = struct {
		CollegeChinese int64 // 大学语文
		ChineseMajor   int64 // 中文专业
		LiteraryExpert int64 // 文学专家
	}{
		CollegeChinese: 41,
		ChineseMajor:   42,
		LiteraryExpert: 43,
	}
)

// Grades 年级列表 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.中考 31.高一 32.高二 33.高考 41.大学语文 42.中文专业 43.文学专家
var Grades = `[
{"id":10,"name":"学龄前","isOpen":true},
{"id":11,"name":"一年级","isOpen":true},
{"id":12,"name":"二年级","isOpen":true},
{"id":13,"name":"三年级","isOpen":true},
{"id":14,"name":"四年级","isOpen":true},
{"id":15,"name":"五年级","isOpen":true},
{"id":16,"name":"六年级","isOpen":true},
{"id":20,"name":"中小衔接","isOpen":true},
{"id":21,"name":"初一","isOpen":true},
{"id":22,"name":"初二","isOpen":true},
{"id":23,"name":"中考","isOpen":true},
{"id":31,"name":"高一","isOpen":true},
{"id":32,"name":"高二","isOpen":true},
{"id":33,"name":"高考","isOpen":true},
{"id":41,"name":"大学语文","isOpen":true},
{"id":42,"name":"中文专业","isOpen":true},
{"id":43,"name":"文学专家","isOpen":true}
]`

// GetGradesList 获取年级列表
func GetGradesList() ([]Grade, error) {
	var grades []Grade
	err := json.Unmarshal([]byte(Grades), &grades)
	if err != nil {
		return nil, errs.WithErr(err)
	}

	return grades, nil
}

// GetUserGradeInfo 获取用户年级信息
func GetUserGradeInfo(grade int64) string {
	grades, _ := GetGradesList()
	for _, v := range grades {
		if grade == v.ID {
			return v.Name
		}
	}
	return ""
}
