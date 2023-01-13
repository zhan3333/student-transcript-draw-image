package transcript

import (
	"fmt"
	"strconv"
	"strings"
)

// Transcript 成绩单
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

// IsGradeValid 成绩是否有效
// 可以接收 [0, 100] 的数字型成绩
// 也可以接收 甲乙丙丁 的文本型成绩
// 允许空成绩
func IsGradeValid(grade string) bool {
	grade = strings.TrimSpace(grade)
	if grade == "" {
		return true
	}
	if isValidTextGrade(grade) {
		return true
	}
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

// ConvertMainGradeToRating 转换成绩为评级
func ConvertMainGradeToRating(grade string) string {
	grade = strings.TrimSpace(grade)
	if grade == "" {
		return ""
	}
	if isValidTextGrade(grade) {
		return grade
	}
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

// ConvertSecondaryGradeToRating 次要科目
func ConvertSecondaryGradeToRating(grade string) string {
	grade = strings.TrimSpace(grade)
	if grade == "" {
		return ""
	}
	if isValidTextGrade(grade) {
		return grade
	}
	g, _ := strconv.ParseFloat(grade, 64)
	if g >= 80 {
		return "甲"
	} else {
		return "乙"
	}
}

var validTextGrades = map[string]bool{
	"甲":  true,
	"乙":  true,
	"丙":  true,
	"丁":  true,
	"缺":  true,
	"甲+": true,
	"甲-": true,
	"乙+": true,
	"乙-": true,
	"丙+": true,
	"丙-": true,
}

func isValidTextGrade(grade string) bool {
	if validTextGrades[grade] {
		return true
	}
	return false
}
