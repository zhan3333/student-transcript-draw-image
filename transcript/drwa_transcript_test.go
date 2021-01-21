package transcript_test

import (
	"fmt"
	"student-scope-send/transcript"
	"testing"
)

func TestNewDrawTranscript(t *testing.T) {
	var err error
	d := transcript.NewDrawTranscript(
		"../0001.jpg",
		"../0001-new.jpg",
		"../fonts/MSYH.TTC",
		//"../fonts/AR-PL-SungtiL-GB.ttf",
		transcript.Transcript{
			Name:           "陈晞文",
			Class:          "二一班",
			Grades:         []string{"甲", "甲", "甲", "甲", "甲", "甲", "甲"},
			StudentComment: "君不见黄河之水天上来，奔流到海不复回。君不见高堂明镜悲白发，朝如青丝暮成雪。",
			ParentComment:  "君不见黄河之水天上来，奔流到海不复回。君不见高堂明镜悲白发，朝如青丝暮成雪。",
			TeacherComment: "君不见黄河之水天上来，奔流到海不复回。君不见高堂明镜悲白发，朝如青丝暮成雪。人生得意须尽欢，莫使金樽空对月。天生我材必有用，千金散尽还复来。烹羊宰牛且为乐，会须一饮三百杯。",
		})
	err = d.ReadTemplate()
	if err != nil {
		fmt.Println(err)
	}
	err = d.Draw()
	if err != nil {
		fmt.Println(err)
	}
	err = d.Save()
	if err != nil {
		fmt.Println(err)
	}
}
