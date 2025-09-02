//go:build js && wasm

package main

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strings"
	"syscall/js"
)

func main() {
	// 初始化随机数种子（TinyGo 兼容）
	rand.Seed(42)
	
	// 注册测试函数
	js.Global().Set("wasmTest", js.FuncOf(testFunction))
	
	// 注册混淆函数
	js.Global().Set("obfuscateJS", js.FuncOf(obfuscateJS))
	
	// 注册验证函数
	js.Global().Set("validateJS", js.FuncOf(validateJS))
	
	// 设置就绪标志
	js.Global().Set("wasmReady", js.ValueOf(true))
	
	// 保持程序运行
	<-make(chan bool)
}

// 测试函数
func testFunction(this js.Value, args []js.Value) interface{} {
	return map[string]interface{}{
		"message": "WASM is working!",
		"success": true,
	}
}

// 混淆配置结构
type ObfuscatorConfig struct {
	IdentifierObfuscation   bool `json:"identifierObfuscation"`
	StringEncryption        bool `json:"stringEncryption"`
	ControlFlowFlattening   bool `json:"controlFlowFlattening"`
	DeadCodeInjection       bool `json:"deadCodeInjection"`
	ExpressionDecomposition bool `json:"expressionDecomposition"`
	CompactCode             bool `json:"compactCode"`
	PreserveComments        bool `json:"preserveComments"`
}

// JavaScript 混淆函数
func obfuscateJS(this js.Value, args []js.Value) interface{} {
	// 添加 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			// TinyGo 兼容：简化错误处理
		}
	}()

	if len(args) < 2 {
		return map[string]interface{}{
			"success": false,
			"error":   "需要提供代码和配置参数",
		}
	}

	code := args[0].String()
	configStr := args[1].String()

	// 解析配置
	var config ObfuscatorConfig
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "配置解析失败: " + err.Error(),
		}
	}

	// 执行混淆
	obfuscatedCode, err := performObfuscationSafe(code, config)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "混淆失败: " + err.Error(),
		}
	}

	// 计算统计信息
	stats := map[string]interface{}{
		"originalSize":   len(code),
		"obfuscatedSize": len(obfuscatedCode),
		"compression":    float64(len(obfuscatedCode)) / float64(len(code)),
	}

	return map[string]interface{}{
		"success": true,
		"code":    obfuscatedCode,
		"stats":   stats,
	}
}

// JavaScript 验证函数
func validateJS(this js.Value, args []js.Value) interface{} {
	// 添加 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			// TinyGo 兼容：简化错误处理
		}
	}()

	if len(args) < 1 {
		return map[string]interface{}{
			"success": false,
			"error":   "需要提供代码参数",
		}
	}

	code := args[0].String()
	
	// 执行验证
	valid, errors := validateJavaScript(code)

	return map[string]interface{}{
		"success": true,
		"valid":   valid,
		"errors":  errors,
	}
}

// 执行实际的混淆操作
func performObfuscation(code string, config ObfuscatorConfig) string {
	result := code
	
	// 移除注释（如果不保留）
	if !config.PreserveComments {
		result = removeComments(result)
	}
	
	// 标识符混淆
	if config.IdentifierObfuscation {
		result = obfuscateIdentifiers(result)
	}
	
	// 字符串加密
	if config.StringEncryption {
		result = encryptStrings(result)
	}
	
	// 控制流平坦化
	if config.ControlFlowFlattening {
		result = flattenControlFlow(result)
	}
	
	// 死代码注入功能已移除
	// 表达式分解功能已移除
	
	// 代码压缩
	if config.CompactCode {
		result = compactCode(result)
	}
	
	return result
}

// 安全的混淆函数，带错误处理
func performObfuscationSafe(code string, config ObfuscatorConfig) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			// TinyGo 兼容：简化错误处理
		}
	}()

	result := code
	
	// 移除注释（如果不保留）
	if !config.PreserveComments {
		result = removeComments(result)
	}
	
	// 标识符混淆
	if config.IdentifierObfuscation {
		result = obfuscateIdentifiers(result)
	}
	
	// 字符串加密
	if config.StringEncryption {
		result = encryptStrings(result)
	}
	
	// 控制流平坦化
	if config.ControlFlowFlattening {
		result = flattenControlFlow(result)
	}
	
	// 死代码注入功能已移除
	// 表达式分解功能已移除
	
	// 代码压缩
	if config.CompactCode {
		result = compactCode(result)
	}
	
	return result, nil
}

// 移除注释
func removeComments(code string) string {
	// 移除单行注释
	singleLineRegex := regexp.MustCompile(`//.*$`)
	code = singleLineRegex.ReplaceAllString(code, "")
	
	// 移除多行注释
	multiLineRegex := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	code = multiLineRegex.ReplaceAllString(code, "")
	
	return code
}

// 彻底移除所有注释 - 用于代码压缩
func removeAllComments(code string) string {
	var result strings.Builder
	inString := false
	inSingleQuote := false
	inDoubleQuote := false
	inSingleLineComment := false
	inMultiLineComment := false
	
	for i := 0; i < len(code); i++ {
		char := code[i]
		
		// 检查是否在字符串中
		if !inSingleLineComment && !inMultiLineComment {
			if char == '\'' && !inDoubleQuote {
				if i == 0 || code[i-1] != '\\' {
					inSingleQuote = !inSingleQuote
					inString = inSingleQuote || inDoubleQuote
				}
			} else if char == '"' && !inSingleQuote {
				if i == 0 || code[i-1] != '\\' {
					inDoubleQuote = !inDoubleQuote
					inString = inSingleQuote || inDoubleQuote
				}
			}
		}
		
		// 如果在字符串中，直接添加字符
		if inString {
			result.WriteByte(char)
			continue
		}
		
		// 检查注释开始
		if !inSingleLineComment && !inMultiLineComment {
			if i < len(code)-1 {
				if char == '/' && code[i+1] == '/' {
					inSingleLineComment = true
					i++ // 跳过下一个字符
					continue
				} else if char == '/' && code[i+1] == '*' {
					inMultiLineComment = true
					i++ // 跳过下一个字符
					continue
				}
			}
		}
		
		// 检查注释结束
		if inSingleLineComment {
			if char == '\n' || char == '\r' {
				inSingleLineComment = false
				result.WriteByte(char) // 保留换行符
			}
			continue
		}
		
		if inMultiLineComment {
			if i < len(code)-1 && char == '*' && code[i+1] == '/' {
				inMultiLineComment = false
				i++ // 跳过下一个字符
			}
			continue
		}
		
		// 如果不在注释中，添加字符
		result.WriteByte(char)
	}
	
	return result.String()
}

// 标识符混淆 - 只混淆用户定义的变量和函数名
func obfuscateIdentifiers(code string) string {
	// JavaScript 保留字和内置对象
	reserved := map[string]bool{
		// 关键字
		"var": true, "let": true, "const": true, "function": true,
		"if": true, "else": true, "for": true, "while": true, "do": true,
		"switch": true, "case": true, "default": true, "break": true, "continue": true,
		"return": true, "try": true, "catch": true, "finally": true, "throw": true,
		"new": true, "this": true, "typeof": true, "instanceof": true, "in": true,
		"class": true, "extends": true, "super": true, "static": true,
		"import": true, "export": true, "from": true, "as": true,
		"async": true, "await": true, "yield": true,
		
		// 字面量
		"undefined": true, "null": true, "true": true, "false": true,
		
		// 全局对象和函数
		"console": true, "window": true, "document": true, "global": true,
		"Array": true, "Object": true, "String": true, "Number": true, "Boolean": true,
		"Date": true, "Math": true, "JSON": true, "RegExp": true, "Error": true,
		"Promise": true, "Symbol": true, "Map": true, "Set": true, "WeakMap": true, "WeakSet": true,
		"parseInt": true, "parseFloat": true, "isNaN": true, "isFinite": true,
		"setTimeout": true, "setInterval": true, "clearTimeout": true, "clearInterval": true,
		"encodeURIComponent": true, "decodeURIComponent": true, "encodeURI": true, "decodeURI": true,
	}
	
	// 如果代码为空，直接返回
	if strings.TrimSpace(code) == "" {
		return code
	}
	
	// 收集用户定义的标识符
	userIdentifiers := make(map[string]bool)
	
	// 1. 收集函数声明
	funcDeclRegex := regexp.MustCompile(`function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`)
	funcMatches := funcDeclRegex.FindAllStringSubmatch(code, -1)
	for _, match := range funcMatches {
		if len(match) > 1 && !reserved[match[1]] {
			userIdentifiers[match[1]] = true
		}
	}
	
	// 2. 收集变量声明
	varDeclRegex := regexp.MustCompile(`(?:var|let|const)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*[=;,]`)
	varMatches := varDeclRegex.FindAllStringSubmatch(code, -1)
	for _, match := range varMatches {
		if len(match) > 1 && !reserved[match[1]] {
			userIdentifiers[match[1]] = true
		}
	}
	
	// 3. 收集函数参数
	funcParamRegex := regexp.MustCompile(`function[^(]*\(([^)]*)\)`)
	paramMatches := funcParamRegex.FindAllStringSubmatch(code, -1)
	for _, match := range paramMatches {
		if len(match) > 1 {
			paramList := match[1]
			paramNames := regexp.MustCompile(`[a-zA-Z_$][a-zA-Z0-9_$]*`).FindAllString(paramList, -1)
			for _, param := range paramNames {
				if !reserved[param] {
					userIdentifiers[param] = true
				}
			}
		}
	}
	
	// 生成混淆映射
	identifierMap := make(map[string]string)
	counter := 0
	for identifier := range userIdentifiers {
		counter++
		identifierMap[identifier] = generateObfuscatedName(counter)
	}
	
	// 如果没有需要混淆的标识符，直接返回原代码
	if len(identifierMap) == 0 {
		return code
	}
	
	// 替换标识符
	result := code
	for original, obfuscated := range identifierMap {
		// 只替换完整的标识符，不替换对象属性
		// 使用更简单的边界匹配，避免负向前瞻
		regex := regexp.MustCompile(`\b` + regexp.QuoteMeta(original) + `\b`)
		
		// 手动检查是否是对象属性调用
		result = regex.ReplaceAllStringFunc(result, func(match string) string {
			// 查找匹配位置
			index := strings.Index(result, match)
			if index >= 0 {
				// 检查后面是否紧跟着点号
				afterMatch := index + len(match)
				if afterMatch < len(result) {
					remaining := strings.TrimLeft(result[afterMatch:], " \t\n\r")
					if strings.HasPrefix(remaining, ".") {
						// 如果后面是点号，不替换
						return match
					}
				}
			}
			return obfuscated
		})
	}
	
	return result
}

// 生成混淆后的标识符名称
func generateObfuscatedName(counter int) string {
	// 简化的混淆策略，兼容 TinyGo
	switch rand.Intn(3) {
	case 0:
		return "_" + intToString(counter)
	case 1:
		return "$_" + intToString(counter)
	default:
		return generateRandomName(6)
	}
}

// 简单的整数转字符串函数
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0'+n%10)}, result...)
		n /= 10
	}
	return string(result)
}

// 生成随机名称
func generateRandomName(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_$"
	result := make([]byte, length)
	
	// 第一个字符不能是数字
	result[0] = chars[rand.Intn(len(chars))]
	
	// 后续字符可以包含数字
	allChars := chars + "0123456789"
	for i := 1; i < length; i++ {
		result[i] = allChars[rand.Intn(len(allChars))]
	}
	
	return string(result)
}

// 字符串加密
func encryptStrings(code string) string {
	// 分别匹配单引号和双引号字符串
	singleQuoteRegex := regexp.MustCompile(`'([^'\\]|\\.)*'`)
	doubleQuoteRegex := regexp.MustCompile(`"([^"\\]|\\.)*"`)
	
	// 处理单引号字符串
	code = singleQuoteRegex.ReplaceAllStringFunc(code, func(match string) string {
		content := match[1 : len(match)-1]
		return encryptString(content)
	})
	
	// 处理双引号字符串
	code = doubleQuoteRegex.ReplaceAllStringFunc(code, func(match string) string {
		content := match[1 : len(match)-1]
		return encryptString(content)
	})
	
	return code
}

// 加密单个字符串
func encryptString(content string) string {
	// 跳过空字符串和很短的字符串
	if len(content) <= 1 {
		return "'" + content + "'"
	}
	
	// 选择加密策略，优先使用更兼容的方法
	strategy := rand.Intn(4)
	
	switch strategy {
	case 0:
		// 字符编码 - 最兼容
		return encodeStringAsCharCodes(content)
	case 1:
		// 十六进制编码 - 兼容性好
		return encodeStringAsHex(content)
	case 2:
		// Unicode 编码
		return encodeStringAsUnicode(content)
	case 3:
		// 简单的字符替换
		return encodeStringAsCharReplace(content)
	default:
		return "'" + content + "'"
	}
}

// 字符编码加密
func encodeStringAsCharCodes(content string) string {
	var parts []string
	for _, char := range content {
		parts = append(parts, "String.fromCharCode("+intToString(int(char))+")")
	}
	return "(" + strings.Join(parts, "+") + ")"
}

// Base64 编码加密
func encodeStringAsBase64(content string) string {
	// 简单的 Base64 编码实现
	encoded := base64Encode(content)
	// 使用自定义的 base64 解码函数，避免依赖 atob
	return "(function(s){var chars='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';var result='';for(var i=0;i<s.length;i+=4){var a=chars.indexOf(s[i])||0;var b=chars.indexOf(s[i+1])||0;var c=chars.indexOf(s[i+2])||0;var d=chars.indexOf(s[i+3])||0;result+=String.fromCharCode((a<<2)|(b>>4));if(s[i+2]!=='=')result+=String.fromCharCode(((b&15)<<4)|(c>>2));if(s[i+3]!=='=')result+=String.fromCharCode(((c&3)<<6)|d);}return result;})('"+encoded+"')"
}

// 十六进制编码加密
func encodeStringAsHex(content string) string {
	var parts []string
	for _, char := range content {
		hex := intToHex(int(char))
		if len(hex) == 1 {
			hex = "0" + hex
		}
		parts = append(parts, "\\x"+hex)
	}
	return "'" + strings.Join(parts, "") + "'"
}

// Unicode 编码加密
func encodeStringAsUnicode(content string) string {
	var parts []string
	for _, char := range content {
		hex := intToHex(int(char))
		for len(hex) < 4 {
			hex = "0" + hex
		}
		parts = append(parts, "\\u"+hex)
	}
	return "'" + strings.Join(parts, "") + "'"
}

// 简单的整数转十六进制函数
func intToHex(n int) string {
	if n == 0 {
		return "0"
	}
	
	hexChars := "0123456789abcdef"
	var result []byte
	for n > 0 {
		result = append([]byte{hexChars[n%16]}, result...)
		n /= 16
	}
	return string(result)
}

// 字符替换加密
func encodeStringAsCharReplace(content string) string {
	var parts []string
	for _, char := range content {
		parts = append(parts, "String.fromCharCode("+intToString(int(char))+")")
	}
	return "(" + strings.Join(parts, "+") + ")"
}

// 简单的 Base64 编码
func base64Encode(input string) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	
	for i := 0; i < len(input); i += 3 {
		var b1, b2, b3 byte
		b1 = input[i]
		if i+1 < len(input) {
			b2 = input[i+1]
		}
		if i+2 < len(input) {
			b3 = input[i+2]
		}
		
		result.WriteByte(chars[b1>>2])
		result.WriteByte(chars[((b1&0x03)<<4)|((b2&0xf0)>>4)])
		
		if i+1 < len(input) {
			result.WriteByte(chars[((b2&0x0f)<<2)|((b3&0xc0)>>6)])
		} else {
			result.WriteByte('=')
		}
		
		if i+2 < len(input) {
			result.WriteByte(chars[b3&0x3f])
		} else {
			result.WriteByte('=')
		}
	}
	
	return result.String()
}

// 控制流平坦化
func flattenControlFlow(code string) string {
	// 简单的控制流混淆
	switchVar := generateRandomName(8)
	
	// 包装在 switch 语句中
	return "\nvar " + switchVar + " = 0;\nwhile (true) {\n\tswitch (" + switchVar + ") {\n\t\tcase 0:\n\t\t\t" + code + "\n\t\t\t" + switchVar + " = 1;\n\t\t\tbreak;\n\t\tcase 1:\n\t\t\treturn;\n\t}\n}"
}

// 死代码注入功能已移除

// 表达式分解功能已移除

// 代码压缩 - 安全版本，保持语法正确性
func compactCode(code string) string {
	// 第一步：彻底移除所有注释
	code = removeAllComments(code)
	
	// 第二步：处理换行和空白字符
	// 将多个连续的空白字符（包括换行）替换为单个空格
	// 但要保留语句分隔符的作用
	
	// 先处理换行符，确保语句正确分隔
	lines := strings.Split(code, "\n")
	var processedLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// 如果行不以分号、大括号、小括号结尾，可能需要添加分号
			if !strings.HasSuffix(line, ";") && 
			   !strings.HasSuffix(line, "{") && 
			   !strings.HasSuffix(line, "}") && 
			   !strings.HasSuffix(line, ",") {
				// 检查是否是控制结构
				trimmed := strings.TrimSpace(line)
				if !strings.HasPrefix(trimmed, "if") && 
				   !strings.HasPrefix(trimmed, "else") && 
				   !strings.HasPrefix(trimmed, "for") && 
				   !strings.HasPrefix(trimmed, "while") && 
				   !strings.HasPrefix(trimmed, "do") && 
				   !strings.HasPrefix(trimmed, "switch") && 
				   !strings.HasPrefix(trimmed, "case") && 
				   !strings.HasPrefix(trimmed, "default") && 
				   !strings.HasPrefix(trimmed, "function") && 
				   !strings.HasPrefix(trimmed, "var") && 
				   !strings.HasPrefix(trimmed, "let") && 
				   !strings.HasPrefix(trimmed, "const") {
					line += ";"
				}
			}
			processedLines = append(processedLines, line)
		}
	}
	
	// 合并所有行
	code = strings.Join(processedLines, " ")
	
	// 第三步：移除多余的空白字符，但保留关键字后的空格
	keywordRegex := regexp.MustCompile(`\b(var|let|const|function|if|else|for|while|do|switch|case|default|break|continue|return|try|catch|finally|throw|new|typeof|instanceof|in|class|extends|super|static|import|export|from|as|async|await|yield)\s+`)
	keywords := keywordRegex.FindAllString(code, -1)
	keywordMap := make(map[string]string)
	for i, keyword := range keywords {
		placeholder := "__KEYWORD_" + intToString(i) + "__"
		keywordMap[placeholder] = keyword
		code = strings.Replace(code, keyword, placeholder, 1)
	}
	
	// 移除多余的空白字符
	spaceRegex := regexp.MustCompile(`\s+`)
	code = spaceRegex.ReplaceAllString(code, " ")
	
	// 恢复关键字
	for placeholder, keyword := range keywordMap {
		code = strings.Replace(code, placeholder, keyword, 1)
	}
	
	// 移除行首行尾空格
	code = strings.TrimSpace(code)
	
	// 第四步：安全地处理操作符，特别注意 ++ 和 -- 操作符
	// 先保护 ++ 和 -- 操作符
	code = strings.ReplaceAll(code, "++", "__PLUSPLUS__")
	code = strings.ReplaceAll(code, "--", "__MINUSMINUS__")
	code = strings.ReplaceAll(code, "+=", "__PLUSEQUAL__")
	code = strings.ReplaceAll(code, "-=", "__MINUSEQUAL__")
	code = strings.ReplaceAll(code, "*=", "__MULTIPLYEQUAL__")
	code = strings.ReplaceAll(code, "/=", "__DIVIDEEQUAL__")
	code = strings.ReplaceAll(code, "==", "__EQUALEQUAL__")
	code = strings.ReplaceAll(code, "===", "__EQUALEQUALEQUAL__")
	code = strings.ReplaceAll(code, "!=", "__NOTEQUAL__")
	code = strings.ReplaceAll(code, "!==", "__NOTEQUALEQUAL__")
	code = strings.ReplaceAll(code, "<=", "__LESSEQUAL__")
	code = strings.ReplaceAll(code, ">=", "__GREATEREQUAL__")
	code = strings.ReplaceAll(code, "&&", "__ANDAND__")
	code = strings.ReplaceAll(code, "||", "__OROR__")
	
	// 移除单个操作符周围的空格
	operatorRegex := regexp.MustCompile(`\s*([+\-*/%=<>!&|^~?:;,(){}[\]])\s*`)
	code = operatorRegex.ReplaceAllStringFunc(code, func(match string) string {
		operator := regexp.MustCompile(`[+\-*/%=<>!&|^~?:;,(){}[\]]`).FindString(match)
		// 对于某些操作符，需要保留前后的空格
		switch operator {
		case "+", "-":
			// 加减号可能是一元操作符，需要小心处理
			return " " + operator + " "
		case "=", "<", ">":
			// 比较和赋值操作符保留空格
			return " " + operator + " "
		default:
			// 其他操作符可以紧贴
			return operator
		}
	})
	
	// 恢复复合操作符
	code = strings.ReplaceAll(code, "__PLUSPLUS__", "++")
	code = strings.ReplaceAll(code, "__MINUSMINUS__", "--")
	code = strings.ReplaceAll(code, "__PLUSEQUAL__", "+=")
	code = strings.ReplaceAll(code, "__MINUSEQUAL__", "-=")
	code = strings.ReplaceAll(code, "__MULTIPLYEQUAL__", "*=")
	code = strings.ReplaceAll(code, "__DIVIDEEQUAL__", "/=")
	code = strings.ReplaceAll(code, "__EQUALEQUAL__", "==")
	code = strings.ReplaceAll(code, "__EQUALEQUALEQUAL__", "===")
	code = strings.ReplaceAll(code, "__NOTEQUAL__", "!=")
	code = strings.ReplaceAll(code, "__NOTEQUALEQUAL__", "!==")
	code = strings.ReplaceAll(code, "__LESSEQUAL__", "<=")
	code = strings.ReplaceAll(code, "__GREATEREQUAL__", ">=")
	code = strings.ReplaceAll(code, "__ANDAND__", "&&")
	code = strings.ReplaceAll(code, "__OROR__", "||")
	
	// 移除不必要的分号前的空格
	semicolonRegex := regexp.MustCompile(`\s*;\s*`)
	code = semicolonRegex.ReplaceAllString(code, ";")
	
	// 移除逗号后的空格，但保留逗号前的空格（可能需要）
	commaRegex := regexp.MustCompile(`,\s+`)
	code = commaRegex.ReplaceAllString(code, ",")
	
	// 移除大括号内外的空格
	braceRegex := regexp.MustCompile(`\s*{\s*`)
	code = braceRegex.ReplaceAllString(code, "{")
	braceRegex2 := regexp.MustCompile(`\s*}\s*`)
	code = braceRegex2.ReplaceAllString(code, "}")
	
	// 移除小括号内外的空格
	parenRegex := regexp.MustCompile(`\s*\(\s*`)
	code = parenRegex.ReplaceAllString(code, "(")
	parenRegex2 := regexp.MustCompile(`\s*\)\s*`)
	code = parenRegex2.ReplaceAllString(code, ")")
	
	return code
}

// 验证 JavaScript 代码
func validateJavaScript(code string) (bool, []string) {
	var errors []string
	
	// 检查代码是否为空
	if len(strings.TrimSpace(code)) == 0 {
		errors = append(errors, "代码不能为空")
		return false, errors
	}
	
	// 检查括号匹配
	if !checkBracketMatching(code) {
		errors = append(errors, "括号不匹配")
	}
	
	// 检查引号匹配
	if !checkQuoteMatching(code) {
		errors = append(errors, "引号不匹配")
	}
	
	// 检查基本语法错误
	if !checkBasicSyntax(code) {
		errors = append(errors, "存在基本语法错误")
	}
	
	return len(errors) == 0, errors
}

// 检查括号匹配
func checkBracketMatching(code string) bool {
	stack := []rune{}
	brackets := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	
	inString := false
	var stringChar rune
	
	for i, char := range code {
		// 处理字符串
		if char == '"' || char == '\'' {
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				// 检查是否被转义
				if i > 0 && rune(code[i-1]) != '\\' {
					inString = false
				}
			}
			continue
		}
		
		if inString {
			continue
		}
		
		// 检查开括号
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		}
		
		// 检查闭括号
		if closing, exists := brackets[char]; exists {
			if len(stack) == 0 {
				return false
			}
			
			if stack[len(stack)-1] != closing {
				return false
			}
			
			stack = stack[:len(stack)-1]
		}
	}
	
	return len(stack) == 0
}

// 检查引号匹配
func checkQuoteMatching(code string) bool {
	inSingleQuote := false
	inDoubleQuote := false
	
	for i, char := range code {
		switch char {
		case '\'':
			if !inDoubleQuote {
				// 检查是否被转义
				if i > 0 && rune(code[i-1]) != '\\' {
					inSingleQuote = !inSingleQuote
				}
			}
		case '"':
			if !inSingleQuote {
				// 检查是否被转义
				if i > 0 && rune(code[i-1]) != '\\' {
					inDoubleQuote = !inDoubleQuote
				}
			}
		}
	}
	
	return !inSingleQuote && !inDoubleQuote
}

// 检查基本语法
func checkBasicSyntax(code string) bool {
	// 检查是否有未闭合的函数
	functionRegex := regexp.MustCompile(`function\s+\w+\s*\([^)]*\)\s*\{`)
	functions := functionRegex.FindAllString(code, -1)
	
	// 简单检查：函数数量应该合理
	if len(functions) > 100 {
		return false
	}
	
	// 检查是否有明显的语法错误模式
	errorPatterns := []string{
		`\}\s*\{`,     // 连续的大括号
		`\)\s*\(`,     // 连续的小括号
		`;;+`,         // 多个分号
		`\+\+\+`,      // 多个加号
		`---`,         // 多个减号
	}
	
	for _, pattern := range errorPatterns {
		if matched, _ := regexp.MatchString(pattern, code); matched {
			return false
		}
	}
	
	return true
}