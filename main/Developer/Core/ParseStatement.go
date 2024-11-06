package Core

import (
	"strings"
)

type Condition struct {
	Key      string
	Operator string
	Value    string
}

type Parser struct{}

// 涉及语法解析器的编写,较难

// convertLogicalExpression 精简化语句,方便进行逆波兰运算 思路为重写方法,不使用正则,而是用一个指针从左往右依次进行扫描并处理
func convertLogicalExpression(expr string) (string, []Condition) { // todo 存在隐形bug: protocol=icmp && banner=icmp_166 无法解析,因为这个指纹非正常预期
	var keyOperatorValue []Condition
	isDoubleQuotationMarkClosure := false
	isBreak := false
	var transformedExpr string
	tempValueStr := ""
	var stack []rune
	runesForExpr := []rune(expr)
	for index := 0; index < len(runesForExpr); index++ {
		if isDoubleQuotationMarkClosure { // 如果进入了",则 keyOperatorValue 开始记录
			tempValueStr = tempValueStr + string(runesForExpr[index])
		}
		isBreak = false
		if index != 0 && runesForExpr[index-1] != '\\' && runesForExpr[index] == '"' {
			if isDoubleQuotationMarkClosure {
				stack = stack[:len(stack)-1]
				isDoubleQuotationMarkClosure = false
				keyOperatorValue[len(keyOperatorValue)-1].Value = tempValueStr // keyOperatorValue 记录完毕
				keyOperatorValue[len(keyOperatorValue)-1].Value = keyOperatorValue[len(keyOperatorValue)-1].Value[:len(keyOperatorValue[len(keyOperatorValue)-1].Value)-1]
				tempValueStr = ""
			} else {
				stack = append(stack, runesForExpr[index]) // 入栈表示进入了"中,transformedExpr 开始不进行记录
				isDoubleQuotationMarkClosure = true
			}
		}

		if len(stack) != 0 || (index+10 > len(runesForExpr)) {
			continue
		}
		switch runesForExpr[index] {
		case ' ':
			continue
		case '(':
			transformedExpr = transformedExpr + "("
			isBreak = true
			break
		case ')':
			transformedExpr = transformedExpr + ")"
			isBreak = true
			break
		case '|':
			transformedExpr = transformedExpr + "|"
			isBreak = true
			break
		case '&':
			transformedExpr = transformedExpr + "&"
			isBreak = true
			break
		}
		if !isBreak {
			switch { // todo bug, cert/server/protocol/banner 在fofa中被特殊使用,不能直接放在header中,否则指纹会失真
			case string(runesForExpr[index:index+5]) == "body=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "=", Value: ""})
				transformedExpr += "body"
				index += 4
				break
			case string(runesForExpr[index:index+6]) == "body!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "!=", Value: ""})
				transformedExpr += "body"
				index += 5
				break
			case string(runesForExpr[index:index+6]) == "title=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "=", Value: ""})
				transformedExpr += "body"
				index += 5
				break
			case string(runesForExpr[index:index+7]) == "title!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "!=", Value: ""})
				transformedExpr += "body"
				index += 6
				break
			case string(runesForExpr[index:index+7]) == "banner=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "=", Value: ""})
				transformedExpr += "body"
				index += 6
				break
			case string(runesForExpr[index:index+8]) == "banner!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "body", Operator: "!=", Value: ""})
				transformedExpr += "body"
				index += 7
				break
			case string(runesForExpr[index:index+7]) == "header=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "=", Value: ""})
				transformedExpr += "header"
				index += 6
				break
			case string(runesForExpr[index:index+8]) == "header!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 7
				break
			case string(runesForExpr[index:index+5]) == "cert=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "=", Value: ""})
				transformedExpr += "header"
				index += 4
				break
			case string(runesForExpr[index:index+6]) == "cert!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 5
				break
			case string(runesForExpr[index:index+9]) == "protocol=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "=", Value: ""})
				transformedExpr += "header"
				index += 8
				break
			case string(runesForExpr[index:index+10]) == "protocol!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 9
				break
			case string(runesForExpr[index:index+7]) == "server=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "=", Value: ""})
				transformedExpr += "header"
				index += 6
				break
			case string(runesForExpr[index:index+8]) == "server!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 7
				break
			case string(runesForExpr[index:index+5]) == "port=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 4
				break
			case string(runesForExpr[index:index+6]) == "port!=":
				keyOperatorValue = append(keyOperatorValue, Condition{Key: "header", Operator: "!=", Value: ""})
				transformedExpr += "header"
				index += 5
				break
			}
		}
	}
	return transformedExpr, keyOperatorValue
}

// getKeyOperatorValue 获取所有的键值对以及算数运算符
func getKeyOperatorValue(conditions []Condition) ([]string, []string, int, int) {
	var bodyArray []string
	var headerArray []string
	bodyCounter := 0
	headerCounter := 0

	for _, cond := range conditions {
		if cond.Key == "body" || cond.Key == "title" || cond.Key == "banner" {
			bodyArray = append(bodyArray, cond.Value)
			bodyCounter++
		} else if cond.Key == "header" || cond.Key == "cert" || cond.Key == "protocol" || cond.Key == "server" || cond.Key == "port" { // todo bug, port/cert/server/protocol/banner 在fofa中被特殊使用,不能直接放在header中,否则指纹会失真
			headerArray = append(headerArray, cond.Value)
			headerCounter++
		}
	}
	return bodyArray, headerArray, bodyCounter, headerCounter
}

// CheckBalanced 用于判断括号是否闭合
func CheckBalanced(s string) bool {
	// 使用切片作为栈
	var stack []rune

	// 创建一个映射，将每个右括号映射到相应的左括号
	matchingBrackets := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		// 如果是左括号，压入栈
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else if char == ')' || char == '}' || char == ']' {
			// 如果是右括号，检查栈是否为空
			if len(stack) == 0 {
				return false // 右括号没有对应的左括号
			}
			// 检查栈顶是否与当前右括号匹配
			top := stack[len(stack)-1]
			if top != matchingBrackets[char] {
				return false // 当前右括号没有对应的左括号
			}
			// 匹配成功，弹出栈顶元素
			stack = stack[:len(stack)-1]
		}
	}

	// 如果栈为空，表示所有括号均已匹配
	return len(stack) != 0
}

// 计算逆波兰表达式的值 [算法可以优化,但是不想弄]
func evaluatePostfix(expression string, bodyArrayByBool []bool, headerArrayByBool []bool) bool {
	var stackByBool []bool  // 只记录结果,是 true 还是 false
	var originStack []uint8 // 原始栈,记录所有数据
	bodyCounter := 0
	headerCounter := 0
	for i := 0; i < len(expression); i++ {
		if expression[i] == 'h' && expression[i+5] == 'r' { // header
			stackByBool = append(stackByBool, headerArrayByBool[headerCounter])
			originStack = append(originStack, expression[i])
			i = i + 5
			headerCounter++
		} else if expression[i] == 'b' && expression[i+3] == 'y' { // body
			stackByBool = append(stackByBool, bodyArrayByBool[bodyCounter])
			originStack = append(originStack, expression[i])
			i = i + 3
			bodyCounter++
		} else if expression[i] == '!' && expression[i+1] == '=' {
			originStack = append(originStack, expression[i])
			i++
		} else if expression[i] == '(' {
			originStack = append(originStack, expression[i])
		} else if expression[i] == ')' {
			for index := len(originStack) - 1; index > 0; index-- {
				if originStack[index] == '(' {
					originStack = originStack[:len(originStack)-1] // 移除栈顶元素
					break
				}
				originStack = originStack[:len(originStack)-1] // 移除栈顶元素
			}
		} else if expression[i] == '&' {
			// 弹出两个操作数进行与运算
			if len(stackByBool) < 2 {
				return false
				//panic("[-] not enough operands") todo 由于技术问题,因此遇到无法解析的指纹直接进行跳过操作
			}
			b := stackByBool[len(stackByBool)-1]
			a := stackByBool[len(stackByBool)-2]
			stackByBool = stackByBool[:len(stackByBool)-2] // 弹出两个操作数
			stackByBool = append(stackByBool, a && b)      // 将结果压入栈
		} else if expression[i] == '|' {
			// 弹出两个操作数进行或运算
			if len(stackByBool) < 2 {
				return false
				//panic("[-] not enough operands") todo 由于技术问题,遇到无法解析的指纹直接进行跳过操作
			}
			b := stackByBool[len(stackByBool)-1]
			a := stackByBool[len(stackByBool)-2]
			stackByBool = stackByBool[:len(stackByBool)-2] // 弹出两个操作数
			stackByBool = append(stackByBool, a || b)      // 将结果压入栈
		}
	}

	if len(stackByBool) != 1 {
		panic("invalid expression")
	}

	return stackByBool[0]
}

// 获取逆波兰表达式
func reversePolishNotation(expression string) string {
	expression = strings.ReplaceAll(expression, " ", "")
	polishNotationByStr := ""
	var originStack []uint8 // 原始栈,记录所有数据
	for i := 0; i < len(expression); i++ {
		if expression[i] == 'h' && expression[i+5] == 'r' { // header
			polishNotationByStr = polishNotationByStr + "header"
			i = i + 5
		} else if expression[i] == 'b' && expression[i+3] == 'y' || (expression[i] == 't' && expression[i+4] == 'e') { // body
			polishNotationByStr = polishNotationByStr + "body"
			i = i + 3
		} else if expression[i] == '(' {
			originStack = append(originStack, expression[i]) // 所有左括号都要入栈
		} else if expression[i] == ')' { // 若是右括号，则栈不断出栈,直到碰到左括号
			for index := len(originStack) - 1; index >= 0; index-- {
				if originStack[index] == '(' {
					originStack = originStack[:len(originStack)-1] // 移除栈顶元素
					break
				}
				polishNotationByStr = polishNotationByStr + string(originStack[len(originStack)-1])
				originStack = originStack[:len(originStack)-1] // 移除栈顶元素
			}
		} else if expression[i] == '&' && expression[i+1] == '&' {
			originStack = append(originStack, expression[i]) // 运算符入栈
			i++
		} else if expression[i] == '|' && expression[i+1] == '|' {
			originStack = append(originStack, expression[i]) // 运算符入栈
			i++
		}
		//else {
		//	Log.Log.Fatal("[-] 你输入的fofa表达式不符合语法,样例:		((body != \"hello\" && header != \"dsd\") && (body!= \"sb\" && header = \"nihao\")) && body!=\"sbb\" || header!=\"ba\"\n")
		//}
	}
	for i := 0; i < len(originStack); i++ {
		polishNotationByStr = polishNotationByStr + string(originStack[i])
	}
	return polishNotationByStr // 示例中正确的逆波兰表达式为:	body header & body header && body header & |
}
