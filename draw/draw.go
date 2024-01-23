package draw

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
	"student-scope-send/constant"
)

// Drawer 绘制成绩单
type Drawer struct {
	TemplateFileName string
	FontFilePath     string
	m                *image.NRGBA
	// 字体尺寸，默认 8
	FontSize float64
	// 行高， 默认 68
	SpaceHeight int
}

func NewDrawer(templateFileName string, fontFilePath string) *Drawer {
	return &Drawer{
		TemplateFileName: templateFileName,
		FontFilePath:     fontFilePath,
		FontSize:         constant.FontSize,
		SpaceHeight:      constant.SpaceHeight,
	}
}

// ReadTemplate 读取模板
func (d *Drawer) ReadTemplate() error {
	src, err := imaging.Open(d.TemplateFileName)
	if err != nil {
		return fmt.Errorf("打开文件失败: %+v\n", err)
	}
	b := src.Bounds()
	d.m = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(d.m, d.m.Bounds(), src, b.Min, draw.Src)
	return nil
}

// Save 保存结果
func (d *Drawer) Save(outputFilePath string) error {
	sp := strings.Split(outputFilePath, "/")
	dirPath := path.Join(sp[0 : len(sp)-1]...)
	err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建文件夹 %s 失败: %+v", sp[0], err)
	}
	err = imaging.Save(d.m, outputFilePath)
	if err != nil {
		return fmt.Errorf("保存图片失败: %+v\n", err)
	}
	return nil
}

func (d *Drawer) Write(text string, rgba color.RGBA, x int, y int) error {
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
	c.SetFontSize(d.FontSize)
	fontFam, err := d.getFontFamily()
	if err != nil {
		return fmt.Errorf("get font family error: %+v\n", err)
	}
	// 设置字体
	c.SetFont(fontFam)
	// 指定位置
	for i, t := range strings.Split(text, "\n") {
		pt := freetype.Pt(x, y+i*d.SpaceHeight)
		_, err = c.DrawString(t, pt)
		if err != nil {
			return fmt.Errorf("draw error: %+v\n", err)
		}
	}

	return nil
}

func (d *Drawer) getFontFamily() (*truetype.Font, error) {
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
