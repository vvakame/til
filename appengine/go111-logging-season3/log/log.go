package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vvakame/sdlog/aelog"
	"github.com/vvakame/sdlog/buildlog"
	"io/ioutil"
	"os"
	"path"
)

var openedFiles []*os.File

func init() {
	var dumpDir func(dirPath string)
	dumpDir = func(dirPath string) {
		fileInfos, err := ioutil.ReadDir(dirPath)
		if err != nil {
			panic(err)
		}
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				fmt.Printf("file: %s\n", dirPath+"/"+fileInfo.Name())
				continue
			}
			fmt.Printf("dir: %s\n", dirPath+"/"+fileInfo.Name())
			dumpDir(dirPath + "/" + fileInfo.Name())
		}
	}
	dumpDir("/var/log")

	fileNames := []string{
		// 普通に出た族
		"/var/log/test1.log",
		// "/var/log/appengine/test2.log",
		// "/var/log/app_engine/app/app-test3.log",
		// "/var/log/app_engine/app/custom_logs/test4.log",
		// "/var/log/test8.txt", // 拡張子 .log じゃなくてもOKだった

		// 駄目族
		// "/var/test5.log", // 怒られが発生した read-only file system
		// "/test6.log", // 怒られが発生した read-only file system
		// "/tmp/test7.log", // これはログとして出力されない(そりゃそう
	}
	for _, fileName := range fileNames {
		dirName := path.Dir(fileName)
		_ = os.MkdirAll(dirName, os.ModePerm)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		openedFiles = append(openedFiles, file)
	}

	aelog.LogWriter = func(ctx context.Context, logEntry *buildlog.LogEntry) {
		b, err := json.Marshal(logEntry)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		for _, f := range openedFiles {
			_, _ = fmt.Fprintln(f, string(b))
			_ = f.Sync()
		}
	}
}

// Criticalf is like Debugf, but at Critical level.
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	aelog.Criticalf(ctx, format, args...)
}

// Debugf formats its arguments according to the format, analogous to fmt.Printf, and records the text as a log message at Debug level.
// The message will be associated with the request linked with the provided context.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	aelog.Debugf(ctx, format, args...)
}

// Errorf is like Debugf, but at Error level.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	aelog.Errorf(ctx, format, args...)
}

// Infof is like Debugf, but at Info level.
func Infof(ctx context.Context, format string, args ...interface{}) {
	aelog.Infof(ctx, format, args...)
}

// Warningf is like Debugf, but at Warning level.
func Warningf(ctx context.Context, format string, args ...interface{}) {
	aelog.Warningf(ctx, format, args...)
}
