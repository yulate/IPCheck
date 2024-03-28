package cmd

import (
	"IPCheck/core"
	"IPCheck/utils"
	"IPCheck/utils/iploc"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "启动分析",
	Long: `启动分析
	分析等级：
		low 	仅扫描80 443端口
		midden 	扫描80, 443, 7000, 8080, 8081, 8443 端口
		high	扫描21, 22, 23, 80, 81, 82, 88, 8000, 8888, 888, 443, 8443, 5000, 7000 端口
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("当前配置参数：")
		path, _ := cmd.Flags().GetString("path")
		level, _ := cmd.Flags().GetString("level")
		threadNum, _ := cmd.Flags().GetString("thread")
		noScan, _ := cmd.Flags().GetBool("noScan")
		fmt.Println("path:", path)
		fmt.Println("level:", level)
		fmt.Println("maxThreads:", threadNum)
		fmt.Println("noScan", noScan)

		// 运行时间计时
		start := time.Now()

		// 文件去重
		utils.DataDeduplication(path)

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 创建一个Scanner来读取文件内容
		scanner := bufio.NewScanner(file)

		// 创建一个等待组，用于等待所有协程完成
		var wg sync.WaitGroup

		// 统计行数
		lineCount := 0

		maxThreads, err := strconv.Atoi(threadNum)
		if err != nil {
			fmt.Println("转换失败：", err)
			return
		}

		// 限制线程数量的信号量
		semaphore := make(chan struct{}, maxThreads) // 设置最大同时运行的协程数量为10

		// 逐行读取文件内容
		for scanner.Scan() {
			lineCount++
			line := scanner.Text()
			wg.Add(1) // 增加等待组的计数

			semaphore <- struct{}{} // 获取一个信号量

			go func(ip string) {
				defer wg.Done()                // 完成时减少等待组的计数
				defer func() { <-semaphore }() // 归还信号量

				operator := iploc.Check(strings.TrimSpace(ip))
				if !noScan {
					if operator != "联通" && operator != "电信" && operator != "移动" {
						ports := core.CheckPort(ip, level)
						if len(ports) > 0 {
							stringNumbers := make([]string, len(ports))

							for i, num := range ports {
								stringNumbers[i] = fmt.Sprintf("%d", num)
							}
							outstr := fmt.Sprintf("IP：%s  归属：%s 开放端口：%s", ip, operator, strings.Join(stringNumbers, ","))
							fmt.Println(outstr)
						}
					}
				} else {
					fmt.Printf("IP：%s  归属：%s \n", ip, operator)
				}

			}(line)
		}

		// 等待所有协程完成
		wg.Wait()

		// 检查Scanner是否出现错误
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("共%v个独立IP地址，总扫描耗时%v", lineCount, elapsed)
	},
}

func init() {
	runCmd.Flags().StringP("path", "p", "./ip.txt", "需检测ip文件位置")
	runCmd.Flags().StringP("level", "l", "low", "检测等级")
	runCmd.Flags().StringP("thread", "t", "10", "扫描线程数")
	runCmd.Flags().BoolP("noScan", "n", false, "是否跳过存活扫描")
	rootCmd.AddCommand(runCmd)
}
