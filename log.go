package tLog

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultFormat = "%d%t [%l]"

var tLogger *TLogger
var tlogConf map[string]map[string]string
var once sync.Once
var format []string
var sol = &stdOutLogger{stdLogger: &stdLogger{level: Info}}

func init() {
	err := loadConf()
	formatAnalysis(viper.GetString("tlog.format"))
	if err != nil {
		sol.warn(fmt.Sprintln("【TLog msg】tLog配置文件加载异常：", *err, "，启用默认配置"), time.Now(), "")
	}

	tlogOut := viper.GetString("tlog.out")

	if tlogOut != "" {
		outs := strings.Split(tlogOut, ",")
		tl := &TLogger{}
		registerOuts(outs, tl)
		tLogger = tl
	} else {
		// 默认的日志输出
		tLogger = &TLogger{
			register: []logger{getStdOutLogger()},
		}
	}
	sol = nil	// help GC，需要手动释放引用吗？？？
}

type TLogger struct {
	register []logger
}

// 加载tlog.yml
func loadConf() *error{
	workDir, _ := os.Getwd()
	viper.SetConfigName("tlog")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.MergeInConfig()
	if err != nil {
		return &err
	}
	return nil
}

// 解析tlog.conf
func getTlogConf() map[string]map[string]string {
	once.Do(func() {
		get := viper.Get("tlog.conf")
		tlogConf = func() map[string]map[string]string {
			var m = make(map[string]map[string]string)
			switch v := get.(type) {
			case []interface{}:
				for _, v1 := range v {
					switch v2 := v1.(type) {
					case map[interface{}]interface{}:
						for k3, v3 := range v2 {
							e := cast.ToString(k3)
							v := cast.ToStringMapString(v3)
							m[e] = v
						}
					}
				}
				return m
			default:
				return m
			}
		}()
	})
	return tlogConf
}

// "%d%t%m [%l] %f"  有意义字母：d t m l f
func formatAnalysis(f string){
	if f == "" {
		f = defaultFormat
	}
	s := strings.Split(f, "%")
	for i, v := range s {
		if i == 0 && v == "" {
			continue
		}
		if v == "" {
			format = append(format, "%")
			continue
		}
		b := []byte(v)
		if b[0] == 'd' || b[0] == 't' || b[0] == 'l' || b[0] == 'm' || b[0] == 'f' {
			format = append(format, string(b[0]))
			if len(b) > 1 {
				format = append(format, string(b[1:]))
			}
		} else {
			format = append(format, "%" + string(b))
		}
	}
}

// 全局默认level
func getDefaultLevel() int{
	dl := viper.GetString("tlog.level.default")
	if dl == "" {
		return Info
	} else {
		return switchLevel(dl)
	}
}

func switchLevel(l string) int {
	switch strings.ToLower(l) {
	case "trace":
		return Trace
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		sol.warn(fmt.Sprintln("【TLog msg】未知的日志输出级别：", l , "，启用默认配置"), time.Now(), "")
		return Info
	}
}

// 注册logger
func registerOuts(outs []string, tl *TLogger) {
	for _, out := range outs {
		switch strings.TrimSpace(out) {
		case confFile:
			tl.register = append(tl.register, getFileOutLogger())
			break
		case confConsole:
			tl.register = append(tl.register, getStdOutLogger())
			break
		default:
			// 对于其他logger待补充
			sol.warn(fmt.Sprintln("【TLog msg】暂不支持logger类型：", out), time.Now(), "")
			break
		}
	}
}

// "%d%t%m [%l] %f"  有意义字母：d t l m f  todo：这里存在bug
func formatHeader(t time.Time, level string) string {
	var s strings.Builder
	var temp string
	for _, f := range format {
		switch f {
		case "d":
			temp += "2006-01-02"
		case "t":
			temp += " 15:04:05"
		case "m":
			temp += ".000"
		case "l":
			if temp != "" {
				s.WriteString(t.Format(temp))
				temp = ""
			}
			s.WriteString(level)
		case "f":
			if temp != "" {
				s.WriteString(t.Format(temp))
				temp = ""
			}
			var ok bool
			_, file, line, ok := runtime.Caller(3)
			if !ok {
				file = "???"
				line = 0
			}
			s.WriteString(file + ":" + strconv.Itoa(line))
		default:
			if temp != "" {
				s.WriteString(t.Format(temp))
				temp = ""
			}
			s.WriteString(f)
		}
	}
	s.WriteString(": ")
	return s.String()
}

// ---------- 入口 -------------
func GetTLog() *TLogger{
	return tLogger
}

func (l *TLogger) Trace(v ...interface{}) {
	t := time.Now()
	var header string
	for _, r := range l.register {
		header = r.trace(fmt.Sprintln(v...), t, header)
	}
}

func (l *TLogger) Debug(v ...interface{}) {
	t := time.Now()
	var header string
	for _, r := range l.register {
		header = r.debug(fmt.Sprintln(v...), t, header)
	}
}

func (l *TLogger) Info(v ...interface{}) {
	t := time.Now()
	s := fmt.Sprintln(v...)
	var header string
	for _, r := range l.register {
		header = r.info(s, t, header)
	}
}

func (l *TLogger) Warn(v ...interface{}) {
	t := time.Now()
	var header string
	for _, r := range l.register {
		header = r.warn(fmt.Sprintln(v...), t, header)
	}
}

func (l *TLogger) Error(v ...interface{}) {
	t := time.Now()
	var header string
	for _, r := range l.register {
		header = r.error(fmt.Sprintln(v...), t, header)
	}
}


