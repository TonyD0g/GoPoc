package Core

import (
	"GoPoc/main/Log"
	"regexp"
	"strings"
)

type Condition struct {
	Key      string
	Operator string
	Value    string
}

type Parser struct{}

// 涉及语法解析器的编写,较难

// convertLogicalExpression 精简化语句,方便进行逆波兰运算
func convertLogicalExpression(expr string) string {
	// 匹配 (body = "value" && header = "value") 这样的结构
	re := regexp.MustCompile(`\(\s*(\w+)\s*=\s*\".*?\"\s*&&\s*(\w+)\s*=\s*\".*?\"\s*\)`)
	transformedExpr := re.ReplaceAllString(expr, "($1 && $2)")

	// 匹配 (body = "value" || header = "value") 这样的结构
	reOr := regexp.MustCompile(`\(\s*(\w+)\s*=\s*\".*?\"\s*\|\|\s*(\w+)\s*=\s*\".*?\"\s*\)`)
	transformedExpr = reOr.ReplaceAllString(transformedExpr, "($1 || $2)")

	// 匹配 body != "value" => body
	reNotEqual := regexp.MustCompile(`(\w+)\s*!=\s*\".*?\"`)
	transformedExpr = reNotEqual.ReplaceAllString(transformedExpr, "$1")

	// 匹配单独的等式，留最后的逻辑变量
	reEqual := regexp.MustCompile(`(\w+)\s*=\s*\".*?\"`)
	transformedExpr = reEqual.ReplaceAllString(transformedExpr, "$1")

	return transformedExpr
}

// getKeyOperatorValue 获取所有的键值对以及算数运算符
func getKeyOperatorValue(expression string) ([]Condition, []string, []string, int, int) {
	var conditions []Condition
	var bodyArray []string
	var headerArray []string
	bodyCounter := 0
	headerCounter := 0

	// 匹配模式：key operator value
	// 支持 ==, !=, =, >, <, >=, <=
	pattern := `(\w+)\s*(!=|==|=|>|<|>=|<=)\s*"([^"]+)"`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(expression, -1)

	for _, match := range matches {
		if len(match) != 4 {
			continue
		}
		conditions = append(conditions, Condition{Key: match[1], Operator: match[2], Value: match[3]})
		if match[1] == "body" {
			bodyArray = append(bodyArray, match[3])
			bodyCounter++
		} else if match[1] == "header" {
			headerArray = append(headerArray, match[3])
			headerCounter++
		}
	}
	return conditions, bodyArray, headerArray, bodyCounter, headerCounter
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
			// todo [优先级最高] 在这修改
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
				panic("[-] not enough operands")
			}
			b := stackByBool[len(stackByBool)-1]
			a := stackByBool[len(stackByBool)-2]
			stackByBool = stackByBool[:len(stackByBool)-2] // 弹出两个操作数
			stackByBool = append(stackByBool, a && b)      // 将结果压入栈
		} else if expression[i] == '|' {
			// 弹出两个操作数进行或运算
			if len(stackByBool) < 2 {
				panic("[-] not enough operands")
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
		} else if expression[i] == 'b' && expression[i+3] == 'y' { // body
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
		} else {
			Log.Log.Fatal("[-] 你输入的fofa表达式不符合语法,样例:		((body != \"hello\" && header != \"dsd\") && (body!= \"sb\" && header = \"nihao\")) && body!=\"sbb\" || header!=\"ba\"\n")
		}
	}
	for i := 0; i < len(originStack); i++ {
		polishNotationByStr = polishNotationByStr + string(originStack[i])
	}
	// 示例中正确的逆波兰表达式为:	body header & body header && body header & |
	return polishNotationByStr
}
