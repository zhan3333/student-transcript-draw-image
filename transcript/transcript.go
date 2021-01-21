package transcript

import (
	"fmt"
	"strconv"
)

//成绩单
type Transcript struct {
	Name           string
	Class          string
	Grades         []string
	StudentComment string
	ParentComment  string
	TeacherComment string
	Email          string
}

type Transcripts []Transcript

//成绩是否有效
func IsGradeValid(grade string) bool {
	i, err := strconv.ParseFloat(grade, 64)
	if err != nil {
		fmt.Printf("%s is not valid grade: %+v\n", grade, err)
		return false
	}
	if i < 0 {
		return false
	}
	if i > 100 {
		return false
	}
	return true
}

//转换成绩为评级
func ConvertMainGradeToRating(grade string) string {
	g, _ := strconv.ParseFloat(grade, 64)
	if g >= 80 {
		return "甲"
	} else if g >= 70 {
		return "乙"
	} else if g >= 60 {
		return "丙"
	} else {
		return "丁"
	}
}

//次要科目
func ConvertSecondaryGradeToRating(grade string) string {
	g, _ := strconv.ParseFloat(grade, 64)
	if g >= 80 {
		return "甲"
	} else {
		return "乙"
	}
}
