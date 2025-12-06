// Package tools 提供关于SQL语句处理的函数
package tools

import (
	"reflect"
	"strings"
	"unicode"
)

// BuildUpdateSet 根据结构体构建SQL语句中set后面的参数值。拼接基本原则：如果结构体字段对应的值为非零值，则进行拼接。PS：此方法在使用时，结构体需带有db的Tag标注。
//
// Params: 1、data：带有db tag的结构体；
//
// 返回结果：
// string：拼接完成的set后面的字符串，带有占位符（%s和？），输出结果示例：name=? , age =?
// []interface{}：拼接字符串对应的参数值的切片。输出示例：[Alice 25]
func BuildUpdateSet(data interface{}) (string, []interface{}) {
	var builder strings.Builder
	var args []interface{}
	reflectType := reflect.TypeOf(data).Elem()
	reflectValue := reflect.ValueOf(data).Elem()
	numField := reflectValue.NumField()
	for i := 0; i < numField; i++ {
		field := reflectValue.Field(i)
		fieldValue := field.Interface()
		if !reflect.DeepEqual(fieldValue, reflect.Zero(field.Type()).Interface()) {
			dbFieldName := reflectType.Field(i).Tag.Get("db")
			// ID不参与update操作
			// if dbFieldName == "id" {
			//	continue
			// }
			builder.WriteString(dbFieldName)
			builder.WriteString("=?,")
			args = append(args, fieldValue)
		}
	}

	// 删除最后一个逗号
	str := strings.TrimRight(builder.String(), ",")

	return str, args
}

// BuildUpdateSetFromMap 根据指定的map来拼装SQL的set部分和对应的参数
func BuildUpdateSetFromMap(updateFields map[string]interface{}) (string, []interface{}) {
	var builder strings.Builder
	var args []interface{}
	keys := reflect.ValueOf(updateFields).MapKeys()
	for i, key := range keys {
		field := key.String()
		value := updateFields[field]
		builder.WriteString(field)
		builder.WriteString("=?")
		args = append(args, value)
		// 如果不是最后一个字段，添加逗号和空格
		if i != len(keys)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String(), args
}

// FilterConditions 生成有效条件的查询条件映射
func FilterConditions(req interface{}) map[string]interface{} {
	conditions := make(map[string]interface{})

	// 使用反射遍历 req 的字段
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// 过滤分页条件
		if field.Name == "PageReq" {
			continue
		}

		// 转换为下划线格式
		fieldName := toSnakeCase(field.Name)

		// 如果字段是嵌套结构体，递归处理
		if field.Type.Kind() == reflect.Struct {
			// 对结构体字段进行递归调用
			subConditions := FilterConditions(value)
			for subKey, subValue := range subConditions {
				// 将子字段展平到外层，并加上前缀
				conditions[subKey] = subValue
			}
		} else {
			// 处理基本字段类型
			if !isZero(value) {
				conditions[fieldName] = value
			}
		}
	}

	return conditions
}

// 判断值是否为零值
func isZero(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case string:
		return v == ""
	case bool:
		return !v
	case float64:
		return v == 0.0
	case nil:
		return true
	default:
		return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
	}
}

// toSnakeCase 将驼峰命名转换为下划线命名
func toSnakeCase(name string) string {
	var result []rune
	for i, char := range name {
		if i > 0 && unicode.IsUpper(char) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
	}
	return string(result)
}
