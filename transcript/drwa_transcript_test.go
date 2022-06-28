package transcript_test

import (
	"fmt"
	"testing"

	"student-scope-send/transcript"
)

func TestNewDrawTranscript(t *testing.T) {
	var err error
	d := transcript.NewDrawTranscript(
		"../0004.jpg",
		"../testdata/0004-out.jpg",
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

func TestMakeNewLine(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		length int
	}{
		{
			name: "s1",
			s: `思维跳跃如脱兔，绿茵赛场似猛虎。
沉淀积累多涵养，动静相宜人人夸。`,
			length: 26,
		}, {
			name:   "s2",
			s:      `你聪明乖巧，非常懂事，对待学习有积极的上进心，希望今后更加专注、踏实，勤字当头。`,
			length: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2 := transcript.MakeNewLine(tt.s, tt.length)
			t.Logf("\n%s", s2)
		})
	}
}
