package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)
import "github.com/360EntSecGroup-Skylar/excelize/v2"

var file string // 文件路径

var excelCmd = &cobra.Command{
	Use:   "excel",
	Short: "表格命令",
}

var convertTranscriptCmd = &cobra.Command{
	Use:   "convert-transcript",
	Short: "将表格中的成绩转换为甲乙丙丁并保存",
	RunE: func(cmd *cobra.Command, args []string) error {
		excel, err := excelize.OpenFile(file)
		if err != nil {
			return fmt.Errorf("open %s: %w", file, err)
		}
		sheet := excel.GetSheetList()[0]
		rows, err := excel.GetRows(sheet)
		if err != nil {
			return fmt.Errorf("get rows: %w", err)
		}
		for i, row := range rows {
			for j, cell := range row {
				if rate, ok := textTToRate(cell); ok {
					axis, err := excelize.CoordinatesToCellName(j+1, i+1)
					if err != nil {
						return fmt.Errorf("get i=%d, j=%d axis: %w", i, j, err)
					}
					if err = excel.SetCellStr(sheet, axis, string(rate)); err != nil {
						return fmt.Errorf("set %s val: %w", axis, err)
					}
				}
			}
		}
		if err := excel.Save(); err != nil {
			return fmt.Errorf("save excel: %w", err)
		}
		fmt.Println("save ok")
		return nil
	},
}

type Rate string

const (
	ARate = "甲"
	BRate = "乙"
	CRate = "丙"
	DRate = "丁"
)

func textTToRate(s string) (Rate, bool) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", false
	}
	if i >= 80 {
		return ARate, true
	} else if i >= 70 {
		return BRate, true
	} else if i >= 60 {
		return CRate, true
	} else {
		return DRate, true
	}
}

func init() {
	excelCmd.AddCommand(convertTranscriptCmd)
	convertTranscriptCmd.Flags().StringVarP(&file, "file", "f", "", "文件路径")
}
