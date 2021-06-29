import React from 'react';
import './App.css';
import {Button, Col, Divider, Layout, message, Row, Space, Upload} from "antd";
import Text from "antd/es/typography/Text";
import {Content, Header} from "antd/es/layout/layout";
import {UploadOutlined} from "@ant-design/icons";

type APPState = {
    taskID: string,
    status: string,
    process: number,
    msg: string,
    loading: boolean,
}

console.log(process.env)

const version = "1.0.0"
const uploadURL = process.env.REACT_APP_SERVER_URL + "/api/upload"
const queryURL = process.env.REACT_APP_SERVER_URL + "/api/query"

export default class App extends React.Component<{}, APPState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            taskID: "",
            status: "",
            process: 0,
            msg: "",
            loading: false,
        }
        this.upload = this.upload.bind(this);
    }

    upload(info: any) {
        console.log("upload", info)
        if (info.file.status !== 'uploading') {
            console.log(info.file, info.fileList);
        }
        if (info.file.status === 'done') {
            message.success(`${info.file.name} file uploaded successfully`);
            this.setState({
                taskID: info.file.response.task_id,
                status: "pending"
            })
        } else if (info.file.status === 'error') {
            this.setState({
                taskID: "",
                status: "failed",
                msg: "上传失败: " + info.file.error
            })
            message.error(`${info.file.name} file upload failed.`);
        }
    }

    enableUpload(): boolean {
        return this.state.taskID === ""
    }

    query() {
        if (this.state.taskID === "") {
            message.warn("无正在执行的任务")
            return
        }
        fetch(queryURL + "?task_id=" + this.state.taskID)
            .then(response => {
                response.json().then((body) => {
                    console.log(response, body)
                    if (response.status === 200) {
                        window.open(body["url"])
                        return
                    } else if (response.status === 202) {
                        this.setState({
                            status: body["status"],
                            process: body["process"],
                            msg: body["msg"]
                        })
                    } else {
                        message.error("查询出错: " + body["msg"])
                    }
                })
            })
    }

    enableQuery(): boolean | undefined {
        if (this.state.taskID === "") {
            return false
        }
        return ["process", "pending", "succeed"].indexOf(this.state.status) !== -1;
    }

    render() {
        return (
            <Layout>
                <Header>
                    <Text className={"header_text"}>华师附属保利南湖小学成绩单系统 {version}</Text>
                </Header>
                <Content>
                    <Row>
                        <Col>
                            <Upload name={"file"} action={uploadURL} onChange={this.upload} accept={".xls,.xlsx"}>
                                <Button icon={<UploadOutlined/>} disabled={!this.enableUpload()}>上传成绩单表格</Button>
                            </Upload>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Space direction={"vertical"}>
                                <Text>TaskID: {this.state.taskID}</Text>
                                <Text>Status: {this.state.status}</Text>
                                <Text>Msg:{this.state.msg}</Text>
                                <Text>Process: {this.state.process} %</Text>
                            </Space>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Button onClick={() => this.query()} disabled={!this.enableQuery()}>查询/下载</Button>
                        </Col>
                    </Row>
                    <Divider/>
                    <Row>
                        <Col>
                            <Space direction={"vertical"}>
                                <Text>Tooltip</Text>
                                <Text>1. 需要重新上传请刷新网页</Text>
                                <Text>2. 刷新网页面前会丢失已上传的文件</Text>
                                <Text>2. 点击上传后，需要手动的点击 查询/下载 按钮来更新状态</Text>
                            </Space>
                        </Col>
                    </Row>
                </Content>
            </Layout>
        );
    }
}