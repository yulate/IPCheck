package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/***
github.com/topxeq/tk 接下来的模块都是为了解决依赖这个工具类导致的软件大小变大问题
*/

// 判断目录下文件是否存在
func IfFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil || os.IsExist(err)
}

func Trim(strA string, cutSetA ...string) string {
	if len(cutSetA) < 1 {
		return strings.TrimSpace(strA)
	}

	return strings.Trim(strA, cutSetA[0])
}

// GenerateErrorString 生成一个出错字符串
func GenerateErrorString(errStrA string) string {
	return "TXERROR:" + errStrA
}

// IsErrorString 判断是否表示出错的字符串
func IsErrorString(errStrA string) bool {
	return StartsWith(errStrA, "TXERROR:")
}

// StartsWith 检查字符串strA开始是否是subStrA
func StartsWith(strA string, subStrA string) bool {
	return strings.HasPrefix(strA, subStrA)
}

func SaveStringListBuffered(strListA []string, fileA string, sepA string) string {
	if strListA == nil {
		return GenerateErrorString("invalid parameter")
	}

	if strListA == nil {
		return GenerateErrorString("empty list")
	}

	lenT := len(strListA)

	var errT error

	file, errT := os.Create(fileA)
	if errT != nil {
		return GenerateErrorString(errT.Error())
	}

	defer file.Close()

	wFile := bufio.NewWriter(file)

	for i := 0; i < lenT; i++ {
		_, errT = wFile.WriteString(strListA[i])
		if errT != nil {
			return GenerateErrorString(errT.Error())
		}

		if i != (lenT - 1) {
			_, errT = wFile.WriteString(sepA)
			if errT != nil {
				return GenerateErrorString(errT.Error())
			}
		}
	}

	wFile.Flush()

	return ""
}

// ReadLineFromBufioReader return result string, error and if reach EOF
func ReadLineFromBufioReader(readerA *bufio.Reader) (string, bool, error) {
	if readerA == nil {
		return "", false, Errf("nil reader")
	}

	strT, errT := readerA.ReadString('\n')

	if errT != nil {
		if errT == io.EOF {
			return strT, true, nil
		}

		return strT, false, errT
	}

	return strT, false, nil

}

// Errf wrap fmt.Errorf function
func Errf(formatA string, argsA ...interface{}) error {
	return fmt.Errorf(formatA, argsA...)
}

// GetErrorString 获取出错字符串中的出错原因部分
func GetErrorString(errStrA string) string {
	if StartsWith(errStrA, "TXERROR:") {
		return errStrA[8:]
	} else {
		return errStrA
	}
}
