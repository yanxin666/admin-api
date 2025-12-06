package excel

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"muse-admin/internal/config"
	"reflect"
	"strings"
	"time"
)

type ExportExcel struct {
	oss      config.Oss // 腾讯云OSS配置
	filepath string     // 文件路径 默认：default
	filename string     // 文件名 默认：当前的毫秒时间戳
	expire   int64      // 生效时间 默认：300秒
}

type OptionsExport func(*ExportExcel)

func NewExCelExport(oss config.Oss, opts ...OptionsExport) *ExportExcel {
	svc := &ExportExcel{
		oss:      oss,
		filepath: "default",
		filename: cast.ToString(util.TimeToTimestampMilli(time.Now())),
		expire:   300, // 默认300秒
	}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

// SetFilepath 选填，设置生成文件的路径
func SetFilepath(path string) OptionsExport {
	return func(e *ExportExcel) {
		e.filepath = path
	}
}

// SetFilename 选填，设置生成文件的名称
func SetFilename(name string) OptionsExport {
	return func(e *ExportExcel) {
		e.filename = name
	}
}

// SetExpire 选填，设置生成文件可下载的有效时间
func SetExpire(expire int64) OptionsExport {
	return func(e *ExportExcel) {
		e.expire = expire
	}
}

// GenExcelFile 导出文件，使用方法：根据header中的key自动匹配struct的字段名，header中的value则是Excel的表头，Excel的内容则是struct的值
// header map[Phone][手机号] 比如：Phone对应的则是 struct的Phone字段，表头为“手机号”，内容为“struct的字段对应的value值”
// data 数据源 []struct{Phone string}
func (e *ExportExcel) GenExcelFile(ctx context.Context, header map[string]string, data interface{}) (string, error) {
	// 使用反射获取data的切片值
	v := reflect.ValueOf(data)

	// 检查 data 是否是切片类型
	if v.Kind() != reflect.Slice {
		return "", fmt.Errorf("expected a slice, got %s", v.Kind())
	}

	// 获取切片的长度
	dataLength := v.Len()
	if dataLength == 0 {
		return "", fmt.Errorf("data slice is empty")
	}

	// 初始化excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 创建Sheet
	sheetName := "sheet1"
	defaultSheet, err := f.NewSheet(sheetName)
	if err != nil {
		logz.Errorf(ctx, "[数据导出]create new sheet failed, err:%v", err)
		return "", err
	}

	// 设置默认Sheet
	f.SetActiveSheet(defaultSheet)

	// 创建一个切片来存储有序的表头
	var sorts []string

	// 获取第一组元素的类型，用于提取字段名
	firstElem := v.Index(0)
	typeOfMember := firstElem.Type()
	for i := 0; i < firstElem.NumField(); i++ {
		fieldName := typeOfMember.Field(i).Name

		// header中不存在的代表不导出，需要过滤
		if _, ok := header[fieldName]; !ok {
			continue
		}

		sorts = append(sorts, fieldName)
	}

	// 生成表头
	var titles []string
	for _, field := range sorts {
		titles = append(titles, header[field])
	}

	// 设置表头到sheet
	err = f.SetSheetRow(sheetName, "A1", &titles)
	if err != nil {
		logz.Errorf(ctx, "[数据导出]set sheet row failed, err:%v", err)
		return "", err
	}
	// 设置所有数据
	for i := 0; i < dataLength; i++ {
		row := make([]interface{}, len(sorts))
		// 获取切片中的每一个元素
		memberValue := v.Index(i)
		for j, field := range sorts {
			// 通过反射获取字段值
			fieldValue := memberValue.FieldByName(field)
			if fieldValue.IsValid() {
				// 获取字段的值
				row[j] = fieldValue.Interface()
			} else {
				// 字段值无效时，设置为空
				row[j] = nil
			}
		}

		// 写入数据
		err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &row)
		if err != nil {
			logz.Errorf(ctx, "[数据导出]set sheet row failed, err:%v", err)
			return "", err
		}
	}

	// 写入Buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		logz.Errorf(ctx, "[数据导出]write to buffer failed, err:%v", err)
	}

	// 生成Excel文件
	r := strings.NewReader(buf.String())

	// 上传腾讯云，生成url
	cos := tencent.NewCos(ctx, tencent.CocConf{
		SecretId:  e.oss.SecretId,
		SecretKey: e.oss.SecretKey,
		Appid:     e.oss.Appid,
		Bucket:    e.oss.Bucket,
		Region:    e.oss.Region,
	})
	return cos.Put(ctx, "admin/"+e.filepath+"/"+e.filename+".xlsx", r, e.expire)
}
