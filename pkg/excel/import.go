package excel

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logc"
	"strings"
)

type ImportExcel struct {
	sheetNum      int
	header        map[string]string // map[表头值][对应字段名] 例：map[手机号][phone]
	isHeaderMatch bool              // 是否效验表头个数
}

type Options func(*ImportExcel)

func NewExCelImport(opts ...Options) *ImportExcel {
	svc := &ImportExcel{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

// WithSheetNum 选填，可以指定读取哪一个sheet，默认为第一个，下标为0
func WithSheetNum(num int) Options {
	return func(e *ImportExcel) {
		e.sheetNum = num
	}
}

// WithHeader 选填，可以指定header头效验，若excel表中没有header中的key匹配的表头，那就直接会报错
func WithHeader(header map[string]string) Options {
	return func(excelImport *ImportExcel) {
		excelImport.header = header
		excelImport.isHeaderMatch = true
	}
}

// GetExcelDataOneSheet 导入文件
// filename string 需要读取的文件内容，指的是上传完到腾讯云后，生成的随机文件名
// uploadFilename string 文件本身名称，指的是用户导入上传文件时，本地存储的源文件名称
func (t *ImportExcel) GetExcelDataOneSheet(ctx context.Context, filename string, uploadFilename ...string) ([]map[string]string, error) {
	// 入参里的filename主要是用作文件内容读取
	f, err := excelize.OpenFile(filename)

	// 若有用户上传的文件名，在读取文件内容异常的时候需要把生成的随机文件名替换为用户上传的，方便理解和报警，属于优化项，顾为选填参数
	if uploadFilename[0] != "" {
		filename = uploadFilename[0]
	}

	if err != nil {
		logc.Errorf(ctx, "解析excel文件报错,filename:%s,错误信息:%v", filename, err)
		return nil, err
	}

	// 获取Sheet的名称
	sheetName := f.GetSheetName(t.sheetNum)

	// 获取所有行的数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		logc.Errorf(ctx, "解析excel文件获取数据报错,filename:%s,错误信息:%v", filename, err)
		return nil, err
	}

	var headers []string
	if len(t.header) > 0 {
		// 第一行的列标题转为对应指定的key
		for _, v := range rows[0] {
			key := t.header[strings.TrimSpace(v)]
			headers = append(headers, key)
		}
	} else {
		headers = rows[0]
	}

	headerLen := len(headers)

	// 自定义表头效验
	if t.isHeaderMatch {
		// 验证个数
		if headerLen < len(t.header) {
			errMsg := fmt.Sprintf(`解析excel文件时，获取数据header头小于要求的header头
		filename：%s
		sheet：%s
		原始header数量：%d
		需要实现header数量：%d`, filename, sheetName, headerLen, len(t.header))
			return nil, errors.New(errMsg)
		}
		// 验证自定义表头是否全都满足
		for k, v := range t.header {
			if !util.IsExist(v, headers) {
				errMsg := fmt.Sprintf(`自定义表头未全部满足
		filename：%s
		sheet：%s
		需要实现的header：%s 不存在`, filename, sheetName, k)
				return nil, errors.New(errMsg)
			}
		}
	}

	// 创建一个切片来存储所有行的数据
	var data []map[string]string

	// 遍历所有行（从第二行开始）
	for _, row := range rows[1:] {
		// 创建一个map来存储当前行的数据
		rowData := make(map[string]string)
		for j, cell := range row {
			// 过滤超过表头的其余数据
			if j+1 > len(headers) {
				continue
			}
			// 过滤表头为空的数据
			if headers[j] == "" {
				continue
			}
			// 将列标题作为键，单元格内容作为值
			rowData[headers[j]] = strings.TrimSpace(cell)
		}
		rowLen := len(row)
		if rowLen != headerLen && rowLen < headerLen {
			for i := 0; i < headerLen-rowLen; i++ {
				rowData[headers[rowLen+i]] = ""
			}
		}
		data = append(data, rowData)
	}
	return data, nil
}
