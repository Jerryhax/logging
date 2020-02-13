/*
@Group : Jerryhax
@Author : Jerryhax
@DateTime : 1/11/20 11:36 AM
@File : log_maintain
@Version : v0.0.1.0
*/

// @Description: 对log file 进行维护：
//1，在每天0点分割前天的日志为单独的文件，并以日期+.log为文件名后缀,然后删除两个月以前的日志文件。
//2，每月1日0点30分，对上个月的日志文件以tar.gz进行打包压缩。
package logging

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Maintain() {
	//定时整理前一天的日志,并删除两个月前的日志
	Infof("[cron job] Every day 00:00 maintain log file.")
	cronJob := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
	_, err := cronJob.AddFunc(CronEveryday, maintain)
	if err != nil {
		log.Fatalln(err)
	}

	//定时归档压缩上个自然月的日志
	Infof("[cron job] Every month 1st 00:30 maintain log file.")
	_, err = cronJob.AddFunc(CronEveryMonth, compress)
	if err != nil {
		log.Fatalln(err)
	}
	cronJob.Start()
}

func maintain() {
	Info("start maintain log file")
	L.Lock()
	defer L.Unlock()

	L.file.Close()

	//archive to yesterday
	destPath := L.filePath + fmt.Sprintf("%s_%s.%s",
		strings.TrimSuffix(L.fileName, "."+LogFileExt),
		time.Now().AddDate(0, 0, -1).Format(TimeFormat),
		LogFileExt)

	//rename the log file.
	err := os.Rename(L.filePath+L.fileName, destPath)
	if err != nil {
		Errorf("occur error while renaming log file %s ,error:%s", destPath, err.Error())
	}
	//reopen a today log file
	L.file, err = MustOpen(L.fileName, L.filePath)
	mw := io.MultiWriter(os.Stdout, L.file)
	L.logger = log.New(mw, DefaultPrefix, log.LstdFlags)

	//delete log files two month ago
	oldDestPath := L.filePath + fmt.Sprintf("%s_%s.%s",
		strings.TrimSuffix(L.fileName, "."+LogFileExt),
		time.Now().AddDate(0, 0, -62).Format(TimeFormat),
		LogFileExt)

	if exists := !CheckNotExist(oldDestPath); exists {
		Infof("try to delete log file %s", oldDestPath)
		err = os.Remove(oldDestPath)
		if err != nil {
			Errorf("occur error while deleting log file %s %s", oldDestPath, err.Error())
		}
	}
	Info("finish maintain log file")
}

// compress log file(*.log) to a tar.gz file (*.log.tar.gz)
func compress() {
	Info("start compress log file")
	lastMonty := time.Now().AddDate(0, -1, 0)
	destPath := L.filePath + fmt.Sprintf("%s_%s.%s",
		strings.TrimSuffix(L.fileName, "."+LogFileExt),
		fmt.Sprintf("%d_%02d", lastMonty.Year(), lastMonty.Month()),
		"tar.gz",
	)

	fw, err := os.Create(destPath)
	if err != nil {
		Fatal(err)
	}
	defer fw.Close()
	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// open dir
	dir, err := os.Open(L.filePath)
	if err != nil {
		Error(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		Error(err)
	}
	if len(files) == 0 {
		Info("There are no log files from a month ago")
		return
	}
	for _, file := range files {
		// 跳过文件夹
		if file.IsDir() {
			continue
		}
		//选取上个自然月的文件
		if !strings.Contains(file.Name(), fmt.Sprintf("%d-%02d", lastMonty.Year(), lastMonty.Month())) {
			continue
		}
		// 文件信息头
		hdr, err := tar.FileInfoHeader(file, "")
		if err != nil {
			log.Println(err)
			continue
		}
		// 写信息头
		err = tw.WriteHeader(hdr)
		if err != nil {
			panic(err)
		}

		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + file.Name())
		if err != nil {
			panic(err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			Errorf("compress log file %s error:%s", destPath, err.Error())
		}

		fr.Close()
	}

	Infof("compress log file %s finish!", destPath)
	return
}
