package read

import (
	"fmt"

	"student-scope-send/transcript"
)
import "github.com/360EntSecGroup-Skylar/excelize/v2"

func Read(path string) (*transcript.Transcripts, error) {
	var ts transcript.Transcripts

	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(f.GetSheetList()[0])
	for i, row := range rows {
		if i < 2 || len(row) < 12 {
			continue
		}
		name := row[1]
		if name == "" {
			continue
		}
		for j, grade := range row[2:7] {
			if !transcript.IsGradeValid(grade) {
				return nil, fmt.Errorf("第 %d 行 第 %d 个成绩填写错误", i+1, j+1)
			}
		}
		grades := []string{
			transcript.ConvertSecondaryGradeToRating(row[2]), //道法
			transcript.ConvertMainGradeToRating(row[3]),      //语文
			transcript.ConvertMainGradeToRating(row[4]),      //数学
			transcript.ConvertSecondaryGradeToRating(row[5]), //英语
			transcript.ConvertSecondaryGradeToRating(row[6]), //体育
			transcript.ConvertSecondaryGradeToRating(row[7]), //艺术
			transcript.ConvertSecondaryGradeToRating(row[8]), //综合
		}
		email := ""
		if len(row) > 12 {
			email = row[12]
		}
		ts = append(ts, transcript.Transcript{
			Name:           name,
			Class:          row[0],
			Grades:         grades,
			StudentComment: row[10],
			ParentComment:  row[11],
			TeacherComment: row[9],
			Email:          email,
		})
	}
	return &ts, nil
}
