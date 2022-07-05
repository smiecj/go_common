// package db db common method: db query condition format, db connector...
package db

import (
	"bytes"
	"fmt"
	"strings"
)

type conditionType string
type conditionMethod string

// join 表方式
type JoinMethod string

const (
	conditionTypeAssert conditionType = "assert"
	conditionTypeOr     conditionType = "or"
	conditionTypeAnd    conditionType = "and"

	conditionMethodLike           conditionMethod = "like"
	conditionMethodEqual          conditionMethod = "equal"
	conditionMethodNotEqual       conditionMethod = "not equal"
	conditionMethodIn             conditionMethod = "in"
	conditionMethodNotIn          conditionMethod = "not in"
	conditionMethodSmaller        conditionMethod = "<"
	conditionMethodBigger         conditionMethod = ">"
	conditionMethodSmallerOrEqual conditionMethod = "<="
	conditionMethodBiggerOrEqual  conditionMethod = ">="

	LeftJoin  JoinMethod = "LEFT JOIN"
	RightJoin JoinMethod = "RIGHT JOIN"
)

var (
	methodToKeywordMap = map[string]string{
		string(conditionMethodLike):           "%",
		string(conditionMethodEqual):          "=",
		string(conditionMethodNotEqual):       "!=",
		string(conditionMethodIn):             "IN",
		string(conditionMethodNotIn):          "NOT IN",
		string(conditionMethodSmaller):        "<",
		string(conditionMethodBigger):         ">",
		string(conditionMethodSmallerOrEqual): "<=",
		string(conditionMethodBiggerOrEqual):  ">=",
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
			// fix-对 value 关键字中有`` -- 作为字段标识, in, not in 这种条件，value 前后不需要加上引号
			condFormatStr := "%s %s '%s'"
			if currentCond.Method == conditionMethodIn || currentCond.Method == conditionMethodNotIn ||
				strings.Contains(currentCond.Value, "`") {
				condFormatStr = "%s %s %s"
			}
			buffer.WriteString(fmt.Sprintf(condFormatStr,
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

// join 条件
type joinCondition struct {
	joinMethod JoinMethod
	space      space
	condition  []string
}

type joinConditionSlice []joinCondition

func (conditionSlice joinConditionSlice) ToSQL() string {
	var conditionBuf bytes.Buffer
	for _, currentCondition := range conditionSlice {
		if conditionBuf.Len() != 0 {
			conditionBuf.WriteString(" ")
		}
		conditionBuf.WriteString(fmt.Sprintf("%s %s ON %s", currentCondition.joinMethod, currentCondition.space.GetSpaceName(),
			buildWhereConditionArr(currentCondition.condition...).ToSQL()))
	}
	return conditionBuf.String()
}

type SearchCondition struct {
	Join  joinConditionSlice
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

type updateCondition struct {
	WhereArr whereArr `json:"where"`
	Limit    int
}

// 获取更新条件中的 limit 部分
func (condition updateCondition) GetLimitCondition() string {
	limitCondition := ""
	if condition.Limit != 0 {
		limitCondition = fmt.Sprintf("LIMIT %d", condition.Limit)
	}
	return limitCondition
}

// 获取更新条件中的 where 条件部分
func (condition updateCondition) GetUpdateCondition() string {
	updateCondition := ""
	if len(condition.WhereArr) != 0 {
		updateCondition = fmt.Sprintf("WHERE %s", condition.WhereArr.ToSQL())
	}
	return updateCondition
}
