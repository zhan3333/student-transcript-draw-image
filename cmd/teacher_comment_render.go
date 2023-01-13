package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"image"
	"math/rand"
	"os"
	"student-scope-send/draw"
	"student-scope-send/transcript"
	"student-scope-send/util"
)

var teacherCommentRenderCmd = &cobra.Command{
	Use:   "teacher_comment_render",
	Short: "生成教师评语图片",
	RunE: func(cmd *cobra.Command, args []string) error {
		for index, comment := range comments {
			// todo 测试只操作第一位同学
			if index > 100 {
				continue
			}
			outFile := fmt.Sprintf("files/out/教师评语/%s.jpg", comment.name)
			if util.IsFileExists(outFile) {
				if err := os.Remove(outFile); err != nil {
					return err
				}
			}
			d := draw.NewDrawer(
				templates[rand.Intn(2)],
				font,
			)
			d.FontSize = 2
			d.SpaceHeight = 28
			if err := d.ReadTemplate(); err != nil {
				return err
			}
			if err := d.Write(
				draw.MakeNewLine(fmt.Sprintf("%s，%s", comment.name, comment.comment), 16),
				transcript.Black,
				45,
				45,
			); err != nil {
				return err
			}
			if err := d.Save(outFile); err != nil {
				return err
			}
		}

		return nil
	},
}

type Drawer struct {
	TemplateImages []image.Image
}

var font = "./fonts/MSYH.TTC"

var templates = []string{
	"assets/teacher_comment.png",
	"assets/teacher_comment2.png",
}

var comments = []struct {
	name    string
	comment string
}{
	{name: "test", comment: "送你一只“模范兔”。你总是积极参加班级各项活动，愿意为班级争光，愿意为班级服务。回首这个特殊的学期，老师多么希望你能继续发挥你的模范带头作用，在学习上也争当表率，拿出责任和担当来吧。"},
}
