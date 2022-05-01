package tLog

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const(
	confFile = "file"
	confPath    = "path"
	confName    = "name"
	confLevel   = "level"
	confMaxSize = "maxSize"
	confMaxNum = "maxNum"
)
// 初始化fileOutLogger
func getFileOutLogger() (fol *fileOutLogger) {
	var path string
	var name string
	var level string
	var maxSize string
	var num string

	conf := getTlogConf()

	if conf != nil {
		m := conf[confFile]
		path = m[confPath]
		name = m[confName]
		level = m[confLevel]
		maxSize = m[confMaxSize]
		num = m[confMaxNum]
	}

	// 获取日志文件路径
	wd, _ := os.Getwd()
	if path == "" {
		path = wd + "/log/"
	} else {
		if !strings.HasPrefix(path, "/") && !strings.Contains(path, ":") {
			path = wd + "/" + path
			if !strings.HasSuffix(path, "/") && !strings.HasSuffix(path, "\\") {
				path = path + "/"
			}
		}
	}

	// 获取日志文件名称
	if name == "" {
		s := strings.Split(wd, string(os.PathSeparator))
		name = s[len(s) - 1]
	}

	if !strings.Contains(name, ".") {
		name = name + ".log"
	}

	// 获取日志level
	var iLevel int
	if level == "" {
		iLevel = getDefaultLevel()
	} else {
		iLevel = switchLevel(level)
	}

	// 封装fileOutLogger
	fol = &fileOutLogger{
		path: path,
		name: name,
		stdLogger: &stdLogger{
			level: iLevel,
		},
	}

	// 设置状态
	if maxSize != "" {
		fol.status = 1
		fol.size = sizeFormat(maxSize)
		fol.num = numFormat(num)
	}

	// 设置out
	file := fol.getOutFile()
	fol.preFile = file
	fol.out = file

	return
}

func sizeFormat(s string) int {
	s2 := s[: len(s)-1]
	n, err := strconv.Atoi(s2)
	if err != nil {
		sol.warn(fmt.Sprintln("【TLog msg】file maxSize 配置异常：", err, "，启用默认配置"), time.Now(), "")
		return 10 << 20
	}
	if n < 1 {
		sol.warn(fmt.Sprintln("【TLog msg】file maxSize 配置异常：", s, "，启用默认配置"), time.Now(), "")
		return 10 << 20
	}

	s3 := strings.ToLower(s[len(s)-1:])
	var i int8
	switch s3 {
	case "k":
		i = 10
	case "m":
		i = 20
	case "g":
		i = 30
	default:
		return 10 << 20
	}
	n = n << i
	return n
}

func numFormat(s string) int{
	n, err := strconv.Atoi(s)
	if err != nil || n <= 1 {
		sol.warn(fmt.Sprintln("【TLog msg】file maxNum 配置异常：", n), time.Now(), "")
		return 0
	}
	return n
}


// ------  fileOutLogger  ------
type fileOutLogger struct {
	*stdLogger

	path string
	name string

	size int		// 每个文件大小
	preSize int		// 存储当前已存储大小
	buf  []byte		// 缓冲区
	lock sync.Mutex	// 锁
	status int8		// 状态，0，无状态即未开启文件分割，1，可写入 2，不可写入
	preFile *os.File	// 当前文件对象
	num int			// 自动清理日志文件数

}

func (fo *fileOutLogger) getOutFile() *os.File {
	s := time.Now().Format("20060102150405")
	if fo.preFile == nil {
		_ = createDir(fo.path, 0666)
		// 如果重启任务，将上次的日志重命名
		_ = os.Rename(fo.path + fo.name, fo.path + fo.name + "_" + s) // 报错不管他
	} else {
		_ = fo.preFile.Close()
		_ = os.Rename(fo.path + fo.name, fo.path + fo.name + "_" + s)
	}

	file, _ := os.OpenFile(fo.path + fo.name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if fo.status == 2 {
		fo.preFile = file
		fo.out = file
		fo.status = 1
		fo.refreshBuf()	// 刷新buf
	}
	// 清理日志文件
	if fo.num > 1 {
		go fo.clearFile()	// 异步减轻生成日志文件效率
	}

	return file
}

func (fo *fileOutLogger) clearFile() {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	f, err := os.Open(fo.path)
	if err != nil {
		fo.warn(fmt.Sprintln("【TLog msg】自动清理日志文件异常：" , err) , time.Now(), "")
		return
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fo.warn(fmt.Sprintln("【TLog msg】自动清理日志文件异常：" , err) , time.Now(), "")
	}
	if len(list) <= fo.num {
		return
	}

	// var s = make([]int, len(list))[0:]
	var s []int
	logName := fo.name + "_"
	for _, file := range list {
		fileName := file.Name()
		if strings.HasPrefix(fileName, logName) {
			suf := fileName[len(logName):]
			temp, err := strconv.Atoi(suf)
			if err != nil {
				continue
			}
			s = append(s, temp)
		}
	}
	if len(s) < fo.num {
		return
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i] - s[j] < 0
	})

	s = s[:len(s) - fo.num + 1]
	for _, v := range s {
		err := os.Remove(fo.path + logName + strconv.Itoa(v))
		if err != nil {
			fo.warn(fmt.Sprintln("【TLog msg】自动清理日志文件异常：" , err) , time.Now(), "")
		}
	}
}

func (fo *fileOutLogger) trace(v string, t time.Time, h string) string  {
	if fo.level <= Trace {
		if h == "" {
			h = formatHeader(t, "trace")
		}
		fo.writeFileLog(h + v)
	}
	return h
}
func (fo *fileOutLogger) debug(v string, t time.Time, h string) string {
	if fo.level <= Debug {
		if h == "" {
			h = formatHeader(t, "debug")
		}
		fo.writeFileLog(h + v)
	}
	return h
}
func (fo *fileOutLogger) info(v string, t time.Time, h string) string {
	if fo.level <= Info {
		if h == "" {
			h = formatHeader(t, "info")
		}
		fo.writeFileLog(h + v)
	}
	return h
}
func (fo *fileOutLogger) warn(v string, t time.Time, h string) string  {
	if fo.level <= Warn {
		if h == "" {
			h = formatHeader(t, "warn")
		}
		fo.writeFileLog(h + v)
	}
	return h
}
func (fo *fileOutLogger) error(v string, t time.Time, h string) string  {
	if fo.level <= Error {
		if h == "" {
			h = formatHeader(t, "error")
		}
		fo.writeFileLog(h + v)
	}
	return h
}

func (fo *fileOutLogger) refreshBuf() {
	fo.writeFileLog("")
}

func (fo *fileOutLogger) writeFileLog(v string) {
	if fo.status == 0 {
		_, _ = fo.out.Write([]byte(v))
	} else {
		fo.lock.Lock()
		defer fo.lock.Unlock()
		if fo.status < 2 {
			if len(fo.buf) > 0 {
				n, _ := fo.out.Write(fo.buf)
				fo.buf = fo.buf[:0]		// 通过这种来清空字节数据但是不清容量
				fo.preSize += n
			}
			if v == "" {	// 刷新buf
				return
			}
			n, _ := fo.out.Write([]byte(v))
			fo.preSize += n
			if fo.preSize >= fo.size {
				fo.status = 2
				go fo.getOutFile()
				fo.preSize = 0
			}
		} else {
			fo.buf = append(fo.buf, v...)
		}
	}
}

