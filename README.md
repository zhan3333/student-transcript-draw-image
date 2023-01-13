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

启动 http 服务

```shell
go run main.go server
```

## 功能

### 调整文字坐标

如果你想测试一下文字坐标对不对，运行下面的测试，并查看 testdata/{out}.jpg 文件观察效果。

`go run test student-scope-send/transcript`

### 导出成绩单

如果你想导出成绩单，那么运行 http server:

`go run main.go server`

并通过这些 api 进行操作:

- POST /api/upload 上传成绩单表格
- GET /api/query 查询导出
- GET /api/send 开始使用邮箱发送成绩