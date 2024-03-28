package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DataDeduplication 文件去重
func DataDeduplication(filePath string) {
	if filePath == "" {
		log.Fatal("没有保存数据，无需进行数据清洗")
	}

	fileT, errT := os.Open(filePath)
	if errT != nil {
		log.Fatal("数据清洗的时候打开文件出错")
	}

	// 保证关闭文件
	defer fileT.Close()

	// 创建读取文件的缓冲式
	readerT := bufio.NewReader(fileT)

	// 记录总行数
	countT := 0

	// 记录总字符数
	totalLenT := 0

	// 循环读取文件的每一行
	for true {
		strT, errT := readerT.ReadString('\n')

		// 如果出现错误中止循环
		if errT != nil {
			if errT == io.EOF {
				strT = strings.TrimRight(strT, "\r\n")
				totalLenT += len(strT)
				countT++
			}
			break
		}

		// 去除行尾可能存在的\r和\n字符（Windows中的文本文件一般的行尾结束符是连续的\r\n两个字符）
		strT = strings.TrimRight(strT, "\r\n")

		// 增加总字符数和总行数
		totalLenT += len(strT)
		countT++
	}

	log.Printf("共%v行，平均每行%v个字符。", countT, totalLenT/countT)
	return

}

// DeduplicationOfLargeFileContent 大文件去重
func DeduplicationOfLargeFileContent(filePath string) {
	var err error

	if filePath == "" {
		log.Fatal("没有保存数据，无需进行数据清洗")
		return
	}

	// 判断文件是否存在
	if !IfFileExists(filePath) {
		fmt.Sprintf("文件 %v 不存在", filePath)
		return
	}

	// limitLineCountT限制每个分块文件的大小（行数）
	// 从命令行参数中可以用-size=100000这样的参数来设置，默认为5000000行
	limitLineCount := 100000

	// 总行数
	lineCount := 0

	// 分块文件数
	fileCount := 0

	// 打开原始文件准备读取
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("数据清洗打开文件时出现错误：", err)
		return
	}

	// 创建一个缓冲式读取器对象
	reader := bufio.NewReader(file)

	// ifEOF 用于判断是否读取到了文件末尾
	ifEOF := false

	// 临时变量，用于存储字符串
	var temps string

	// 反复循环从源文件中读取行，直至读到文件末尾
	// 每次读取最多limitLineCount行，写入临时文件中，超出则继续写到下一个临时文件中
	// 临时文件名按数字进行排序，存于变量fileCount中
	for !ifEOF {
		// 分配指定大小的切片，准备放置读取到的文本行
		buf := make([]string, 0, limitLineCount)

		fileCount++

		// 临时文件名将依次为sub00000001.txt、sub00000002.txt...
		subFileName := fmt.Sprintf("sub%08d.txt", fileCount)

		// 默认将临时文件放在执行时的当前目录下
		subPath := filepath.Join("./result/", subFileName)

		// 循环读取limitLineCount次，试图读取limitLineCount行文本
		for j := 0; j < limitLineCount; j++ {
			str, err := reader.ReadString('\n')
			if err != nil {
				// 读到文件结尾时的处理
				if err == io.EOF {
					temps = Trim(str)
					// 判断最后一行是否存在内容，如果存在则将其添加入buf
					if temps != "" {
						buf = append(buf, temps)
					}
					ifEOF = true
				} else {
					log.Fatal("文件读取失败：", err)
					file.Close()
					os.Exit(1)
				}
				break
			}

			temps = Trim(str)

			// 将空行丢弃
			if temps != "" {
				buf = append(buf, temps)
			}
		}

		// 对读取到的最多limitLineCount行文本进行排序
		fmt.Sprintf("正在排序第%d组数据", fileCount)
		sort.Sort(sort.StringSlice(buf))

		// 保存排序后的文本到临时文件
		fmt.Sprintf("正在保存第%d组数据到临时文件%v", fileCount, subPath)
		res := SaveStringListBuffered(buf, subPath, "\n")
		if IsErrorString(res) {
			fmt.Sprintf("保存临时文件%v失败：%v", subPath, GetErrorString(res))
			file.Close()
			os.Exit(1)
		}

		// 记录总共处理的行数
		lineCount += len(buf)
	}

	file.Close()

	fmt.Sprintf("共读取了%v行，写入了%v个临时文件", lineCount, fileCount)

	// 排序写
	fmt.Println("进行多文件排序并去除重复行......")

	// 存放临时文件读取器的变量
	files := make([]*os.File, fileCount)
	readers := make([]*bufio.Reader, fileCount)

	// 用于进行对多个文件读取的第一行进行大小对比排序的变量
	strBuf := make([]string, fileCount)
	compareBuf := make([]int, fileCount)
	selIndex := 0

	// 用于保存当前写入的行，用于去重重复行
	currentLine := ""

	// 统计整体读取的行数和写入的行数
	readCount := 0
	writeCount := 0

	// 打开多个临时文件用于同时读取
	for i := 1; i <= fileCount; i++ {
		subPath := filepath.Join("./result/", fmt.Sprintf("sub%08d.txt", i))

		fmt.Sprintf("打开临时文件%v准备读取", subPath)
		files[i-1], err = os.Open(subPath)
		if err != nil {
			fmt.Sprintf("打开临时文件%v出现错误：%v", subPath, err)
			os.Exit(1)
		}

		readers[i-1] = bufio.NewReader(files[i-1])
	}

	// 创建一个新文件用于写入最终结果
	outputPath := strings.Trim(filePath, ".txt")
	outputFile, err := os.Create("./" + outputPath + "_deduplication.txt")
	fmt.Println(outputPath)
	if err != nil {
		log.Fatal("创建输出文件时发生错误：", err)
		os.Exit(1)
	}

	// 创建写入器
	outputWriter := bufio.NewWriter(outputFile)

	// 用于判断是否是写入的第一行
	// 如果不是第一行m，将再写入每一行文本前先写入一个回车换行符
	notFirstFlag := false

	// 记录一个被关闭了多少个临时文件，表示已经有多少个临时文件被读取完毕
	var closedFile = 0

	// 循环读取并写入结果文件
	for true {
		var line string

		// 是否读取到文件结尾
		var eof bool

		// 从各个文件中都读取一行，空行将被丢弃
		for k := 0; k < fileCount; k++ {
			if readers[k] == nil {
				closedFile++
				continue
			}

			// 如果某个文件对应的一行已空，则再读一行
			if strBuf[k] == "" {
				found := false
				for readers[k] != nil {
					line, eof, err = ReadLineFromBufioReader(readers[k])

					if err != nil {
						fmt.Sprintf("从临时文件%v中读取数据时发生错误：%v", k, err)
						os.Exit(1)
					}

					line = Trim(line)

					if eof {
						readers[k] = nil
						files[k].Close()
					}

					if line == "" {
						continue
					}

					found = true
					break
				}

				if found {
					strBuf[k] = line
				}
			}
		}

		// 进行计数式比对，找出排名最考前的一行
		var compare int

		for ii := 0; ii < fileCount; ii++ {
			compareBuf[ii] = 0
		}

		for ii := 0; ii < (fileCount - 1); ii++ {
			if strBuf[ii] == "" {
				continue
			}

			for jj := ii + 1; jj < fileCount; jj++ {
				if strBuf[jj] == "" {
					compareBuf[ii]++
					continue
				}

				compare = strings.Compare(strBuf[ii], strBuf[jj])
				if compare > 0 {
					compareBuf[jj]++
				} else if compare < 0 {
					compareBuf[ii]++
				}
			}
		}

		max := 0
		for kk := 0; kk < fileCount; kk++ {
			if compareBuf[kk] > max {
				max = compareBuf[kk]
				selIndex = kk
			}
		}

		// 处理只有一个文件时的对比
		if fileCount == 1 && strBuf[0] != "" {
			max = 1
			selIndex = 0
		}

		// 如果所有行都是空行，说明已经读取完毕所有文件，将退出循环
		if max <= 0 {
			fmt.Println("读取缓冲区全部为空")
			break
		}

		readCount++

		// 如果将要写入的一行与上一行一样，说明是重复行，则丢弃
		// 由此实现去除重复行的效果
		// 注意此方法仅对排序后的文本才是正确的
		if currentLine != "" {
			if strBuf[selIndex] == currentLine {
				strBuf[selIndex] = ""
				continue
			}
		}

		currentLine = strBuf[selIndex]
		strBuf[selIndex] = ""

		if notFirstFlag {
			outputWriter.WriteString("\r\n")
		} else {
			notFirstFlag = true
		}

		// 将最终选出的文本行写入结果文件
		_, err = outputWriter.WriteString(currentLine)
		if err != nil {
			log.Fatal("去重结果写入文件错误：", err)
			os.Exit(1)
		}

		writeCount++

		// 所有文件如果都已关闭，说明都已读取完，循环将终止
		if closedFile == fileCount {
			break
		}
	}
	// 由于使用的是bufio，即缓冲方式写入文件，注意一定要用Flush来保证在内存中的数据被确保真正写入文件中
	outputWriter.Flush()
	outputFile.Close()

	// 删除临时文件
	for l := 1; l <= closedFile; l++ {
		temporaryFiles := "./result/" + fmt.Sprintf("sub%08d.txt", l)
		// 先进行判断文件是否存在
		if IfFileExists(temporaryFiles) {
			err := os.Remove(temporaryFiles)
			if err != nil {
				fmt.Println("临时文件删除出现错误：", err)
			}
		}
	}

	fmt.Sprintf("去重操作处理完毕，共写入%v行", writeCount)
}
