package transcript_test

import (
	"github.com/stretchr/testify/assert"
	"student-scope-send/draw"
	"testing"

	"student-scope-send/transcript"
)

func TestNewDrawTranscript(t *testing.T) {
	drawer := draw.NewDrawer(
		"../0005.jpg",
		"../fonts/MSYH.TTC",
		//"../fonts/AR-PL-SungtiL-GB.ttf",
	)
	transcript2 := &transcript.Transcript{
		Name:           "陈晞文",
		Class:          "二一班",
		Grades:         []string{"甲", "甲", "甲", "甲", "甲", "甲", "甲"},
		StudentComment: "君不见黄河之水天上来，奔流到海不复回。君不见高堂明镜悲白发，朝如青丝暮成雪。",
		ParentComment:  "你平时性格内向，不善多言，忠厚老实，跟同学能友好相处，平时能关心集体，值日工作负责，喜爱体育活动，不过，你在学习上还要努力些，作业时更要细心，把字写好。",
		TeacherComment: "君不见黄河之水天上来，奔流到海不复回。君不见高堂明镜悲白发，朝如青丝暮成雪。人生得意须尽欢，莫使金樽空对月。天生我材必有用，千金散尽还复来。烹羊宰牛且为乐，会须一饮三百杯。",
	}
	out := "../testdata/0005-out.jpg"
	if assert.NoError(t, transcript.Draw(drawer, transcript2, out)) {
		t.Logf("save to %s", out)
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
			s2 := draw.MakeNewLine(tt.s, tt.length)
			t.Logf("\n%s", s2)
		})
	}
}
