import React from 'react';
import './App.css';
import {Button, Col, Layout, message, Row, Upload} from "antd";
import Text from "antd/es/typography/Text";
import {Content} from "antd/es/layout/layout";
import {UploadOutlined} from "@ant-design/icons";

type APPState = {
    taskID: string,
    status: string,
    process: number,
    msg: string
}

console.log(process.env)

const uploadURL = process.env.REACT_APP_SERVER_URL + "/upload"
const queryURL = process.env.REACT_APP_SERVER_URL + "/query"

export default class App extends React.Component<{}, APPState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            taskID: "",
            status: "",
            process: 0,
            msg: ""
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

    query() {
        if (this.state.taskID === "") {
            message.warn("无正在执行的任务")
            return
        }
        fetch(queryURL + "?task_id=" + this.state.taskID)
            .then(response => {
                response.json().then((body) => {
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
                <Content>
                    <Row>
                        <Col>
                            <Upload name={"file"} action={uploadURL} onChange={this.upload} accept={".xls,.xlsx"}>
                                <Button icon={<UploadOutlined/>}>Click to Upload</Button>
                            </Upload>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Text>{this.state.taskID}: {this.state.status}: {this.state.msg}: {this.state.process / 100} %</Text>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Button onClick={() => this.query()} disabled={!this.enableQuery()}>查询/下载</Button>
                        </Col>
                    </Row>
                </Content>
            </Layout>
        );
    }
}