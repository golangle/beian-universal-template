package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

// 数据和模板文件的路径
var filePaths []string = []string{"conf/hosts.txt", "template/template.tmpl"}
var configPath string = "conf/default.conf"

// 记录文件的最后修改时间
var lastFileTimes []time.Time

// 应用启动时间
var startTime time.Time = time.Now()

// 记录首次访问的IP地址
var firstIp string = ""

// 刷新间隔，单位秒
var refreshInterval int = 5 //默认5秒

var ip_filter bool = false

// 主机信息结构体
type HostInfo struct {
	Title     string
	BeiAn     string
	CopyRight string
}

// 主机名到数据行的映射
var HostLineMap map[string]HostInfo = map[string]HostInfo{}

// 模板实例
var templateInstance *template.Template

/*
=========================================================================================================================

	加载配置文件

=========================================================================================================================
*/
func loadConfigFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if len(line) == 0 {
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if string(strings.TrimSpace(line)[0]) == "#" {
			continue
		}
		lineSegment := strings.Split(line, ":")
		if len(lineSegment) > 1 {
			name := strings.TrimSpace(lineSegment[0])
			value := strings.TrimSpace(lineSegment[1])
			// 跳过无名配置项
			if len(name) == 0 {
				continue
			}
			// 跳过无值配置项
			if len(value) == 0 {
				continue
			}
			if name == "refresh_interval" {
				if iv, err := strconv.Atoi(value); err == nil {
					refreshInterval = iv
				} else {
					log.Printf("invalid refresh_interval %q: %v", value, err)
				}
			}
			if name == "ip_filter" {
				if value == "true" {
					ip_filter = true
				} else {
					ip_filter = false
				}
			}
		}
	}
}

/*
=========================================================================================================================

	加载数据文件和模板文件

=========================================================================================================================
*/
func loadData() {

	// 热加载时，清空原有数据
	if len(HostLineMap) > 0 {
		lastFileTimes = []time.Time{}
		HostLineMap = map[string]HostInfo{}
	}
	// 记录本次加载文件的修改时间
	for _, pv := range filePaths {
		initFileLastModTime(pv)
	}

	// 解析模板文件
	tmpl, err := template.ParseFiles(filePaths[1])
	if err != nil {
		panic(err)
	}
	templateInstance = tmpl

	file, err := os.Open(filePaths[0]) // 打开文件
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if len(line) == 0 {
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if string(strings.TrimSpace(line)[0]) == "#" {
			continue
		}
		lineSegment := strings.Split(line, "|") //分割标题、主机名、版权信息和备案号

		// 有一个 | 分割符，表示有主机名
		if len(lineSegment) > 1 {
			hosts := strings.Split(lineSegment[1], ",") //分割主机名
			for i := range hosts {
				host := strings.TrimSpace(hosts[i])
				if len(lineSegment) == 4 {
					HostLineMap[host] = HostInfo{
						Title:     strings.TrimSpace(lineSegment[0]), //分配标题
						CopyRight: strings.TrimSpace(lineSegment[2]), //分配版权信息
						BeiAn:     strings.TrimSpace(lineSegment[3]), //分配备案号
					}
				}
				if len(lineSegment) == 3 {
					HostLineMap[host] = HostInfo{
						Title:     strings.TrimSpace(lineSegment[0]), //分配标题
						CopyRight: strings.TrimSpace(lineSegment[2]), //分配版权信息
					}
				}
				if len(lineSegment) == 2 {
					HostLineMap[host] = HostInfo{
						Title: strings.TrimSpace(lineSegment[0]), //分配标题
					}
				}
			}
		} else {
			//如果只写了数据，没有 | 分割符，则此数据赋给备案号。匹配所有主机。
			HostLineMap["*"] = HostInfo{
				BeiAn: strings.TrimSpace(lineSegment[0]), //分配给备案号
			}
		}
	}
}

/*
=========================================================================================================================

	启动时初始化文件时间

===========================================================================================================================
*/

func initFileLastModTime(filePath string) {
	// 首次获取文件修改时间
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lastModTime := fileInfo.ModTime()
	lastFileTimes = append(lastFileTimes, lastModTime)
}

/*
=========================================================================================================================

	检查文件在应用启动后是否被修改,如果被修改则返回true，否则返回false。

===========================================================================================================
*/
func IsNeedReloadData() bool {
	// 获取文件修改时间并比较
	fileInfo, err := os.Stat(filePaths[0]) //配置文件
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	configCurrentModTime := fileInfo.ModTime()

	fileInfo, err = os.Stat(filePaths[1]) // 模板文件
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	templateCurrentModTime := fileInfo.ModTime()

	if configCurrentModTime.After(lastFileTimes[0]) || templateCurrentModTime.After(lastFileTimes[1]) {
		return true
	}
	return false
}

// =========================================================================================================================
//
//	访问日志写入功能
//
// =========================================================================================================================
// 配置文件路径
var configFile = "log/access.txt"

// 文件同步锁
var lock sync.RWMutex

// 写入日志
func writeLog(msg string) {
	lock.Lock()
	defer lock.Unlock()
	//读取日志文件
	fileHandle, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 使用 NewWriterSize()方法定义缓冲区
	buf := bufio.NewWriterSize(fileHandle, len(msg))

	//写入消息
	buf.WriteString(msg)

	//写入磁盘
	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}
}
