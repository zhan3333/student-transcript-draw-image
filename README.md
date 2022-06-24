# 学生成绩图片绘制程序

## todo list

- [x] 更换成绩单背景
- [x] 提供接口
    - [x] 上传成绩单 excel 接口 /api/upload
    - [x] 下载成绩单图压缩包 /api/query
    - [x] 发送邮件 /api/send
- [x] 提供 web 界面调用上述两个接口
- [] 提供登录接口，使用 cookie session 保存用户身份
- [] 提供在线预览界面
- [] 提供全部发送邮件界面

## 命令

转换表格中数字成绩为甲乙丙丁评级

```shell
go run main.go excel convert-transcript -f="/Users/zhan/Downloads/期末成绩单-demo(1).xlsx"
```
