package utils

import (
	"container/list"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/kdada/tinygo/log"
)

// Logger 控制台日志器
var Logger log.Logger
var logger log.Logger

func init() {
	log.Register("f", FriendlyLoggerCreator)
	var err error
	logger, err = log.NewLogger("f", "3")
	if err != nil {
		panic(err)
	}
	Logger, err = log.NewLogger("f", "2")
	if err != nil {
		panic(err)
	}
}

// Debug 输出debug日志
func Debug(info ...interface{}) {
	logger.Debug(info)
}

type FriendlyLogger struct {
	logLevel  log.LogLevel
	logList   *list.List
	logmu     *sync.Mutex
	logWriter log.LogWriter
	closed    bool
	async     bool
	skip      int
}

// NewFriendlyLogger 创建日志记录器
func NewFriendlyLogger(logWriter log.LogWriter, skip int) *FriendlyLogger {
	return &FriendlyLogger{
		log.LogLevelDebug | log.LogLevelInfo | log.LogLevelWarn | log.LogLevelError | log.LogLevelFatal,
		list.New(),
		new(sync.Mutex),
		logWriter,
		false,
		false,
		skip,
	}
}

// writeLog 写入日志
func (this *FriendlyLogger) writeLog(info string, level log.LogLevel) {
	if !this.closed && level&this.logLevel > 0 {
		info = time.Now().Format("2006-01-02 15:04:05.000000 ") + info
		this.logmu.Lock()
		if this.async {
			this.logList.PushBack(info)
		} else {
			this.logWriter.Write(info)
		}
		this.logmu.Unlock()
	}
}

// file 返回调用Logger的函数位置
func (this *FriendlyLogger) position(prefix string) string {
	if this.skip >= 2 {
		var _, file, line, ok = runtime.Caller(this.skip)
		if ok {
			return prefix + file + "(" + strconv.Itoa(line) + "):"
		}
	}
	return prefix
}

// Debug 写入调试信息
func (this *FriendlyLogger) Debug(info ...interface{}) {
	this.writeLog(this.position("[Debug]")+fmt.Sprintf("%+v", info), log.LogLevelDebug)
}

// Info 写入一般信息
func (this *FriendlyLogger) Info(info ...interface{}) {
	this.writeLog(this.position("[Info]")+fmt.Sprintf("%+v", info), log.LogLevelInfo)
}

// Warn 写入警告信息
func (this *FriendlyLogger) Warn(info ...interface{}) {
	this.writeLog(this.position("[Warn]")+fmt.Sprintf("%+v", info), log.LogLevelWarn)
}

// Error 写入错误信息
func (this *FriendlyLogger) Error(info ...interface{}) {
	this.writeLog(this.position("[Error]")+fmt.Sprintf("%+v", info), log.LogLevelError)
}

// Fatal 写入崩溃信息
func (this *FriendlyLogger) Fatal(info ...interface{}) {
	this.writeLog(this.position("[Fatal]")+fmt.Sprintf("%+v", info), log.LogLevelFatal)
}

// LogLevel 得到日志等级是否输出
func (this *FriendlyLogger) LogLevelOutput(level log.LogLevel) bool {
	return this.logLevel&level > 0
}

// SetLogLevel 设置某个日志等级是否输出
func (this *FriendlyLogger) SetLogLevelOutput(level log.LogLevel, output bool) {
	if output {
		this.logLevel |= level
	} else {
		this.logLevel &= ^level
	}
}

// Async 是否异步输出
func (this *FriendlyLogger) Async() bool {
	return this.async
}

// SetAsync 设置是否异步输出
func (this *FriendlyLogger) SetAsync(async bool) {
	var oldAsync = this.async
	this.async = async
	this.logWriter.SetAsync(this.async, this.logList, this.logmu)
	if oldAsync && !this.async && this.logList.Len() > 0 {
		//从异步切换回同步,将尚未异步输出的日志转换为同步输出
		var start *list.Element
		var length = 0
		this.logmu.Lock()
		if this.logList.Len() > 0 {
			start = this.logList.Front()
			length = this.logList.Len()
			this.logList.Init()
		}
		for i := 0; i < length; i++ {
			var v, ok = start.Value.(string)
			if ok {
				this.logWriter.Write(v)
			}
			start = start.Next()
		}
		this.logmu.Unlock()
	}
}

// SetSkip skip为跳过的Caller数量,skip小于2时关闭文件位置记录的功能
func (this *FriendlyLogger) SetSkip(skip int) {
	this.skip = skip
}

// Closed 日志是否已关闭
func (this *FriendlyLogger) Closed() bool {
	return this.closed
}

// Close 关闭日志 关闭后无法再使用
func (this *FriendlyLogger) Close() {
	if !this.Closed() {
		this.closed = true
		this.logWriter.Close()
	}
}

func FriendlyLoggerCreator(param string) (log.Logger, error) {
	var skip, err = strconv.Atoi(param)
	if err != nil {
		fmt.Println(err)
		skip = 2
	}
	return NewFriendlyLogger(log.NewSimpleLogWriter(log.NewConsoleWriter()), skip), nil
}
