package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"student-scope-send/read"
	"student-scope-send/transcript"
	"time"
)

func UploadTranscript(c *gin.Context) {
	var (
		taskID = GenerateTaskID()
	)
	// 单文件
	file, _ := c.FormFile("file")

	// 上传文件至指定目录
	savePath := filepath.Join("files/upload", randomFileName(file.Filename, taskID))
	err := os.MkdirAll(strings.ReplaceAll(savePath, filepath.Base(savePath), ""), os.ModeDir|os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "上传成功",
		"task_id": taskID,
	})
	fmt.Printf("[%s] upload %s success\n", taskID, file.Filename)
	task := Task{
		ID:       taskID,
		FilePath: savePath,
		Status:   "pending",
		Process:  0,
		StartAt:  nil,
		EndAt:    nil,
	}
	if err = task.Cache(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Errorf("缓存失败 %w", err),
		})
		return
	}
	go operator(taskID)
}

func DownloadTranscriptImg(c *gin.Context) {
	task := Task{
		ID: c.Query("task_id"),
	}
	if err := task.Resume(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": fmt.Sprintf("%s 不存在", task.ID),
		})
		return
	}
	var (
		status  int
		msg     string
		process = task.Process
	)
	switch task.Status {
	case "pending":
		status = http.StatusAccepted
		msg = "等待处理中"
	case "process":
		status = http.StatusAccepted
		msg = "处理中"
	case "succeed":
		status = http.StatusOK
		msg = "处理成功"
	case "failed":
		status = http.StatusInternalServerError
		msg = fmt.Sprintf("处理失败: %s", task.Err)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "未知的 task status " + task.Status,
		})
		return
	}

	// 打包下载文件
	if task.Status == "succeed" {
		err := os.MkdirAll("files/export", os.ModeDir|os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		zipPath := path.Join("files", "export", fmt.Sprintf("%s.zip", task.ID))
		if !isFileExists(zipPath) {
			err = Zip(path.Join("files/out", task.ID), zipPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"url": fmt.Sprintf("%s/api/export/%s", os.Getenv("ORIGIN"), filepath.Base(zipPath)),
		})
		return
	}

	c.JSON(status, gin.H{
		"msg":     msg,
		"process": process,
	})
}

type Task struct {
	ID       string     `json:"id"`
	FilePath string     `json:"file_path"`
	Status   string     `json:"status"`
	Process  int        `json:"process"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
	Err      string     `json:"err"`
}

func (t *Task) Cache() error {
	b, _ := json.Marshal(t)
	return RDB.Set(context.Background(), t.Key(), string(b), 0).Err()
}

func (t *Task) Resume() error {
	s, err := RDB.Get(context.Background(), t.Key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s 不存在", t.Key())
		}
		return err
	}
	return json.Unmarshal([]byte(s), t)
}

func (t *Task) SetErr(err error) {
	t.Err = fmt.Sprintf("[%s] %+v", t.Key(), err)
	t.Status = "failed"
	_ = t.Cache()
}

func (t *Task) Key() string {
	return fmt.Sprintf("task_%s", t.ID)
}

func randomFileName(filename string, taskID string) string {
	return fmt.Sprintf("%s_%s", taskID, filepath.Base(filename))
}

const randStrs = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateTaskID() string {
	rand.Seed(time.Now().Unix())
	ret := ""
	for i := 0; i <= 32; i++ {
		ret += string(randStrs[rand.Intn(len(randStrs))])
	}
	return ret
}

func operator(taskID string) {
	var (
		outDir               = path.Join("files", "out", taskID)
		failedSendEmailNames []string
	)
	task := Task{ID: taskID}
	if err := task.Resume(); err != nil {
		err = fmt.Errorf("读取缓存失败: %w\n", err)
		fmt.Printf("%+v\n", err)
		task.SetErr(err)
		return
	}
	task.Status = "process"
	now := time.Now()
	task.StartAt = &now
	ts, err := read.Read(task.FilePath)
	if err != nil {
		err = fmt.Errorf("读取成绩失败: %w\n", err)
		fmt.Printf("%+v\n", err)
		task.SetErr(err)
		return
	}
	fmt.Printf("读取到%d条成绩单数据\n\n\n", len(*ts))
	for i, t := range *ts {
		fmt.Printf("开始处理第 %d 条数据: %s\n", i+1, t.Name)
		outFile := fmt.Sprintf("%s/%s.jpg", outDir, t.Name)
		d := transcript.NewDrawTranscript(
			templateFilePath,
			outFile,
			fontFilePath,
			t,
		)
		if isFileExists(outFile) {
			fmt.Printf("成绩单 %s 已经存在, 如需重新生成, 请先删除对应成绩单文件\n", outFile)
		} else {
			// 生成成绩单
			err = d.ReadTemplate()
			if err != nil {
				fmt.Printf("读取第 %d 个学生 %s 时读取模板失败: %+v\n", i+1, t.Name, err)
				continue
			}
			err = d.Draw()
			if err != nil {
				fmt.Printf("绘制第 %d 个学生 %s 成绩单失败: %+v\n", i+1, t.Name, err)
				continue
			}
			err = d.Save()
			if err != nil {
				fmt.Printf("保存第 %d 个学生 %s 成绩单失败: %+v\n", i+1, t.Name, err)
				continue
			}
		}

		//fmt.Println("开始发送邮件")
		//if isSendEmail(t.Name) {
		//	fmt.Println("缓存发现已经发送过邮件, 跳过发送")
		//} else {
		//	time.Sleep(5 * time.Second)
		//	err = sendEmail(t, outFile)
		//	if err != nil {
		//		fmt.Printf("邮件发送失败: %+v\n", err)
		//		failedSendEmailNames = append(failedSendEmailNames, t.Name)
		//	} else {
		//		setSendEmail(t.Name)
		//		fmt.Println("邮件发送完毕")
		//	}
		//}

		task.Process = 100 * (i + 1) / len(*ts)
		_ = task.Cache()
		fmt.Printf("第 %d 条数据: %s 处理完成\n", i+1, t.Name)
		fmt.Println("---------------------------------")
	}
	task.Process = 100
	task.Status = "succeed"
	now = time.Now()
	task.EndAt = &now
	_ = task.Cache()
	fmt.Printf("发送邮件失败的学生名单: %v\n", failedSendEmailNames)
	fmt.Println("程序执行完毕")
}

func Zip(src_dir string, zip_file_name string) error {
	fmt.Println(src_dir, zip_file_name)

	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, err := os.Create(zip_file_name)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	return filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}
