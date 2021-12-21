// package db db common method: db query condition format, db connector...
package db

import (
	"bytes"
	"fmt"
)

type conditionType string
type conditionMethod string

const (
	conditionTypeAssert conditionType = "assert"
	conditionTypeOr     conditionType = "or"
	conditionTypeAnd    conditionType = "and"

	conditionMethodLike    conditionMethod = "like"
	conditionMethodEqual   conditionMethod = "equal"
	conditionMethodSmaller conditionMethod = "<"
	conditionMethodBigger  conditionMethod = ">"
)

var (
	methodToKeywordMap = map[string]string{
		string(conditionMethodLike):    "%",
		string(conditionMethodEqual):   "=",
		string(conditionMethodSmaller): "<",
		string(conditionMethodBigger):  ">",
	}
	keyWordToMethodMap = map[string]conditionMethod{
		"%": conditionMethodLike,
		"=": conditionMethodEqual,
		"<": conditionMethodSmaller,
		">": conditionMethodBigger,
	}
)

// db where condition json format
// example: cond_example.json
// 参数校验: key、value 需要进行校验，可放在接口层，防止出现 1=1 这种绝对正确的条件
type whereCondition struct {
	Type   conditionType   `json:"type"`
	Key    string          `json:"key"`
	Method conditionMethod `json:"method"`
	Value  string          `json:"value"`
}

// 多个where 条件定义
type whereArr []whereCondition

// where 组合条件 转换成 SQL 查询语句
func (arr whereArr) ToSQL() string {
	buffer := new(bytes.Buffer)
	for _, currentCond := range arr {
		switch currentCond.Type {
		case conditionTypeAssert:
			buffer.WriteString(fmt.Sprintf("%s %s '%s'",
				currentCond.Key, methodToKeywordMap[string(currentCond.Method)], currentCond.Value))
		case conditionTypeAnd, conditionTypeOr:
			buffer.WriteString(fmt.Sprintf(" %s ", currentCond.Type))
		}
	}
	return buffer.String()
}

// 通过传入的条件 生成 whereCondition
// 格式: "name", "equal" / "=", "xiaoming", "and", "grade", "equal" / "=", "3"
func buildWhereConditionArr(args ...string) whereArr {
	retArr := make(whereArr, 0)
	index := 0
	for index < len(args) {
		currentCondition := whereCondition{}
		currentArg := args[index]
		switch currentArg {
		case string(conditionTypeAnd), string(conditionTypeOr):
			currentCondition.Type = conditionType(currentArg)
		default:
			if index+2 >= len(args) || (methodToKeywordMap[args[index+1]] == "" && keyWordToMethodMap[args[index+1]] == "") {
				break
			}
			conditionMethod := conditionMethod(args[index+1])
			if methodToKeywordMap[args[index+1]] == "" {
				conditionMethod = keyWordToMethodMap[args[index+1]]
			}

			currentCondition.Key, currentCondition.Method, currentCondition.Value =
				currentArg, conditionMethod, args[index+2]
			currentCondition.Type = conditionTypeAssert
		}
		retArr = append(retArr, currentCondition)
		index++
	}
	return retArr
}

type SearchCondition struct {
	Order struct {
		Field string `json:"field"`
		Sc    string `json:"sc"`
	} `json:"order"`
	Page struct {
		No    int `json:"no"`
		Limit int `json:"limit"`
	} `json:"page"`
	WhereArr whereArr `json:"where"`
}

type UpdateCondition struct {
	WhereArr whereArr `json:"where"`
	Limit    int
}
