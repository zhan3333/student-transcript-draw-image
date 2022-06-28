package controller

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/olekukonko/tablewriter"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"student-scope-send/app"
	"student-scope-send/read"
	"student-scope-send/transcript"
	"student-scope-send/util"
)

var templateFilePath = "./0004.jpg"
var fontFilePath = "./fonts/MSYH.TTC"

func Upload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	taskID := GenerateTaskID(strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)))
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

func Query(c *gin.Context) {
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
		status     int
		msg        string
		process    = task.Process
		url        string
		table      [][]string
		mailState  string
		mailErrMsg string
	)

	table = append(table, []string{"姓名", "邮箱", "成绩单", "状态", "错误"})
	for _, mail := range task.Mails {
		mailErrMsg = ""
		mailState = "待发送"
		if mail.Error != "" {
			mailErrMsg = mail.Error
			mailState = "失败"
		}
		table = append(table, []string{mail.Name, mail.Email, mail.FilePath, mailState, mailErrMsg})
	}
	for _, mail := range task.SentMails {
		table = append(table, []string{mail.Name, mail.Email, mail.FilePath, "成功", ""})
	}

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
		if !util.IsFileExists(zipPath) {
			err = Zip(path.Join("files/out", task.ID), zipPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
		}
		url = fmt.Sprintf("%s/api/export/%s", os.Getenv("ORIGIN"), filepath.Base(zipPath))
	}

	c.JSON(status, gin.H{
		"url":        url,
		"msg":        msg,
		"process":    process,
		"status":     task.Status,
		"mails":      len(task.Mails),
		"sent_mails": len(task.SentMails),
		"table":      table,
	})
}

func Send(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("task_id 必须传入"),
		})
		return
	}
	if err := send(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("发送邮件发生错误 %+v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

type Task struct {
	ID        string     `json:"id"`
	FilePath  string     `json:"file_path"`
	Status    string     `json:"status"`
	Process   int        `json:"process"`
	StartAt   *time.Time `json:"start_at"`
	EndAt     *time.Time `json:"end_at"`
	Err       string     `json:"err"`
	Mails     []TaskMail `json:"mails"`
	SentMails []TaskMail `json:"sent_mails"`
}

type TaskMail struct {
	Name      string `json:"name"`
	FilePath  string `json:"file_path"`
	Email     string `json:"email"`
	FailCount int    `json:"fail_count"`
	Error     string `json:"-"`
}

func (t *Task) Cache() error {
	b, _ := json.Marshal(t)
	return app.GetRedis().Set(context.Background(), t.Key(), string(b), 24*30*time.Hour).Err()
}

func (t *Task) Resume() error {
	s, err := app.GetRedis().Get(context.Background(), t.Key()).Result()
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

func GenerateTaskID(prefix string) string {
	return prefix + "-" + time.Now().Format("2006010215405")
}

func operator(taskID string) {
	var (
		outDir = path.Join("files", "out", taskID)
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
		if util.IsFileExists(outFile) {
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
		task.Mails = append(task.Mails, TaskMail{Name: d.Transcript.Name, FilePath: d.OutFilePath, Email: d.Transcript.Email})
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
	fmt.Println("程序执行完毕")
}

func send(taskID string) error {
	task := Task{ID: taskID}
	if err := task.Resume(); err != nil {
		err = fmt.Errorf("读取缓存失败: %w\n", err)
		return err
	}
	if len(task.Mails) == 0 {
		return fmt.Errorf("无待发送的邮件")
	}
	var (
		sent []TaskMail
		fail []TaskMail
	)
	fmt.Println("开始发送邮件")
	for _, mail := range task.Mails {
		err := func() error {
			if mail.Email == "" {
				return fmt.Errorf("%s 未配置邮箱", mail.Name)
			}
			if !util.IsFileExists(mail.FilePath) {
				return fmt.Errorf("%s 文件不存在", mail.FilePath)
			}

			err := SendEmail(mail.Email, mail.Name, mail.FilePath)
			if err != nil {
				return fmt.Errorf("邮件发送失败: %w", err)
			}
			return nil
		}()
		if err != nil {
			fmt.Printf("%s %s 发送失败: %+v\n", mail.Name, mail.Email, err)
			mail.FailCount++
			mail.Error = err.Error()
			fail = append(fail, mail)
		} else {
			fmt.Printf("%s %s 发送成功\n", mail.Name, mail.Email)
			sent = append(sent, mail)
		}
		time.Sleep(5 * time.Second)
	}
	task.Mails = fail
	task.SentMails = sent
	fmt.Printf("发送邮件结束, 成功 %d / 失败 %d\n", len(task.SentMails), len(task.Mails))
	if len(task.Mails) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"姓名", "邮箱", "失败次数", "最后一次失败原因"})
		for _, mail := range task.Mails {
			table.Append([]string{mail.Name, mail.Email, string(rune(mail.FailCount)), mail.Error})
		}
		table.Render()
	}
	if err := task.Cache(); err != nil {
		b, _ := json.Marshal(task)
		fmt.Println("cache failed", string(b))
		return err
	}
	return nil
}

func Zip(src_dir string, zip_file_name string) error {
	fmt.Println(src_dir, zip_file_name)

	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, err := os.Create(zip_file_name)
	if err != nil {
		return err
	}
	defer func() { _ = zipfile.Close() }()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer func() { _ = archive.Close() }()

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
			defer func() { _ = file.Close() }()
			_, _ = io.Copy(writer, file)
		}
		return nil
	})
}
