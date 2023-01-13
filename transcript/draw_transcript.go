package transcript

import (
	"fmt"
	"student-scope-send/draw"
)

// Draw 绘制图片并保存结果
func Draw(drawer *draw.Drawer, transcript *Transcript, outputFilePath string) error {
	var err error

	// 读取模版
	// todo 优化一下不需要每次都读取模版
	if err := drawer.ReadTemplate(); err != nil {
		return err
	}

	// 学生姓名
	err = drawer.Write(transcript.Name, Black, studentNameCoordinate.X, studentNameCoordinate.Y)
	if err != nil {
		return err
	}

	// 班级
	err = drawer.Write(transcript.Class, Black, classCoordinate.X, classCoordinate.Y)
	if err != nil {
		return err
	}

	// 成绩
	for i, grade := range transcript.Grades {
		x := gradesCoordinate.X + i*gradesOffset // 初始坐标 + 偏移量
		y := gradesCoordinate.Y
		err := drawer.Write(grade, Black, x, y)
		if err != nil {
			return fmt.Errorf("写第 %d 个成绩失败: %+v\n", i, err)
		}
	}

	// 童言
	err = drawer.Write(draw.MakeNewLine(transcript.StudentComment, 9), Black, studentCommentCoordinate.X, studentCommentCoordinate.Y)
	if err != nil {
		return err
	}

	// 家长
	err = drawer.Write(draw.MakeNewLine(transcript.ParentComment, 11), Black, parentCommentCoordinate.X, parentCommentCoordinate.Y)
	if err != nil {
		return err
	}

	// 教师
	err = drawer.Write(draw.MakeNewLine(transcript.TeacherComment, 26), Black, teacherCommentCoordinate.X, teacherCommentCoordinate.Y)
	if err != nil {
		return err
	}

	// 保存
	if err := drawer.Save(outputFilePath); err != nil {
		return err
	}

	return nil
}
