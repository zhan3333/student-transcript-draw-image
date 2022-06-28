package transcript

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//绘制成绩单
type DrawTranscript struct {
	TemplateFileName string
	OutFilePath      string
	FontFilePath     string
	Transcript       Transcript
	m                *image.NRGBA
}

func NewDrawTranscript(templateFileName string, outFileName string, fontFilePath string, transcript Transcript) *DrawTranscript {
	return &DrawTranscript{TemplateFileName: templateFileName, OutFilePath: outFileName, FontFilePath: fontFilePath, Transcript: transcript}
}

//读取模板
func (d *DrawTranscript) ReadTemplate() error {
	src, err := imaging.Open(d.TemplateFileName)
	if err != nil {
		return fmt.Errorf("打开文件失败: %+v\n", err)
	}
	b := src.Bounds()
	d.m = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(d.m, d.m.Bounds(), src, b.Min, draw.Src)
	return nil
}

//保存结果
func (d *DrawTranscript) Save() error {
	t := strings.Split(d.OutFilePath, "/")
	dirPath := path.Join(t[0 : len(t)-1]...)
	fmt.Println(t, dirPath)
	err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建文件夹 %s 失败: %+v", t[0], err)
	}
	err = imaging.Save(d.m, d.OutFilePath)
	if err != nil {
		return fmt.Errorf("保存图片失败: %+v\n", err)
	}
	return nil
}

//绘制
func (d *DrawTranscript) Draw() error {
	var err error
	// 学生姓名
	err = d.writeName(d.Transcript.Name)
	if err != nil {
		return err
	}

	// 班级
	err = d.writeClass(d.Transcript.Class)
	if err != nil {
		return err
	}

	// 成绩
	err = d.writeGrades(d.Transcript.Grades)
	if err != nil {
		return err
	}

	// 童言
	err = d.writeStudentComment(d.Transcript.StudentComment)
	if err != nil {
		return err
	}

	// 家长
	err = d.writeParentComment(d.Transcript.ParentComment)
	if err != nil {
		return err
	}

	// 教师
	err = d.writeTeacherComment(d.Transcript.TeacherComment)
	if err != nil {
		return err
	}
	return nil
}

//写教师评语, 26个字符一行
func (d *DrawTranscript) writeTeacherComment(comment string) error {
	var x = 780
	var y = 1925
	return d.write(MakeNewLine(comment, 26), black, x, y)
}

// MakeNewLine 根据长度加入换行
func MakeNewLine(s string, length int) string {
	var arr []string
	for _, l := range strings.Split(s, "\n") {
		s2 := ""
		for j, c := range []rune(l) {
			if (j+1)%length == 0 {
				s2 += "\n"
			}
			s2 += string(c)
		}
		arr = append(arr, s2)
	}

	return strings.Join(arr, "\n")
}

// 童言妙语
func (d *DrawTranscript) writeStudentComment(comment string) error {
	var x = 890
	var y = 2655
	return d.write(MakeNewLine(comment, 9), black, x, y)
}

// 家长心语
func (d *DrawTranscript) writeParentComment(comment string) error {
	var x = 2000
	var y = 2610
	return d.write(MakeNewLine(comment, 9), black, x, y)
}

// 成绩
func (d *DrawTranscript) writeGrades(grades []string) error {
	for i, grade := range grades {
		x := 1255 + i*238 // 初始坐标 + 偏移量
		y := 3660
		err := d.write(grade, black, x, y)
		if err != nil {
			return fmt.Errorf("写第 %d 个成绩失败: %+v\n", i, err)
		}
	}
	return nil
}

// 学生姓名
func (d *DrawTranscript) writeName(name string) error {
	var x = 1400
	var y = 900
	return d.write(name, black, x, y)
}

// 班级
func (d *DrawTranscript) writeClass(class string) error {
	var x = 2273
	var y = 900
	return d.write(class, black, x, y)
}

func (d *DrawTranscript) write(text string, rgba color.RGBA, x int, y int) error {
	// todo 要实现字体间距
	//spacing := 1.5 // 字间距
	c := freetype.NewContext()

	c.SetDPI(512)
	c.SetClip(d.m.Bounds())
	c.SetDst(d.m)
	c.SetHinting(font.HintingFull)

	// 设置文字颜色AdobeSongStd.otf
	c.SetSrc(image.NewUniform(rgba))
	// 设置字体大小
	c.SetFontSize(8)
	fontFam, err := d.getFontFamily()
	if err != nil {
		return fmt.Errorf("get font family error: %+v\n", err)
	}
	// 设置字体
	c.SetFont(fontFam)
	// 指定位置
	for i, t := range strings.Split(text, "\n") {
		pt := freetype.Pt(x, y+i*68)
		_, err = c.DrawString(t, pt)
		if err != nil {
			return fmt.Errorf("draw error: %+v\n", err)
		}
	}

	return nil
}

func (d *DrawTranscript) getFontFamily() (*truetype.Font, error) {
	// 这里需要读取中文字体，否则中文文字会变成方格
	fontBytes, err := ioutil.ReadFile(d.FontFilePath)
	if err != nil {
		fmt.Println("read file error:", err)
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("parse font error:", err)
		return &truetype.Font{}, err
	}

	return f, err
}
