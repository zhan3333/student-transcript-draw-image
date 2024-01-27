package cmd

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"student-scope-send/app"
	"student-scope-send/controller"
	"student-scope-send/util"
	"time"
)

var send bool
var emailFile string
var reportCardFolder string
var cacheSendStudentNameKey = "compare_and_send:send"

func init() {
	compareAndSendCmd.Flags().BoolVar(&send, "send", false, "send email")
	compareAndSendCmd.Flags().StringVar(&emailFile, "email_file", "", "email excel")
	compareAndSendCmd.Flags().StringVar(&reportCardFolder, "report_card_folder", "", "report card folder")
	compareAndSendCmd.MarkFlagsRequiredTogether("email_file", "report_card_folder")
	rootCmd.AddCommand(compareAndSendCmd)
}

var compareAndSendCmd = &cobra.Command{
	Use:   "compare_and_send",
	Short: "compare email and report cards, then send emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		students, err := readStudents(emailFile)
		if err != nil {
			return fmt.Errorf("read students: %w", err)
		}
		//for _, student := range students {
		//	fmt.Println(student.Name, student.Email)
		//}

		cards, err := readReportCards(reportCardFolder)
		if err != nil {
			return fmt.Errorf("read report cards: %w", err)
		}
		//for _, card := range cards {
		//	fmt.Printf("%+v\n", card)
		//}

		matches, err := match(students, cards)
		if err != nil {
			return fmt.Errorf("match: %w", err)
		}
		if send {
			var sendFailed []*Match
			for _, match := range matches {
				fmt.Printf("%s,%s,%s start sent email\n", match.Student.Name, match.Student.Email, match.ReportCard.Name)

				// determine whether the email was sent
				sent, err := isSent(match.Student.Email)
				if err != nil {
					return fmt.Errorf("determine whether the email was sent faild")
				}
				if sent {
					fmt.Printf("%s alread sent, skip\n", match.Student.Name)
					continue
				}

				// sent email
				err = controller.SendEmail(match.Student.Email, match.Student.Name, match.ReportCard.Path)
				if err != nil {
					sendFailed = append(sendFailed, match)
					fmt.Printf("sent email to %s failed: %s\n", match.Student.Name, err)
				} else {
					fmt.Printf("sent email to %s success\n", match.Student.Name)
					if err := setSent(match.Student.Email); err != nil {
						fmt.Printf("set sent failed: %s\n", err)
					}
				}
				fmt.Println("sleep 5s...")
				time.Sleep(5 * time.Second)
			}
			if len(sendFailed) > 0 {
				fmt.Printf("%d student email send failed\n", len(sendFailed))
			} else {
				fmt.Println("all email send")
			}
		}

		return nil
	},
}

type Student struct {
	Name  string
	Email string
}

type Match struct {
	Student    *Student
	ReportCard *ReportCard
}

func isSent(studentName string) (bool, error) {
	return app.GetRedis().SIsMember(context.Background(), cacheSendStudentNameKey, studentName).Result()
}

func setSent(studentName string) error {
	return app.GetRedis().SAdd(context.Background(), cacheSendStudentNameKey, studentName).Err()
}

func match(students []*Student, cards []*ReportCard) ([]*Match, error) {
	if len(students) == 0 {
		return nil, fmt.Errorf("students number=0")
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("report cards number=0")
	}

	var matches []*Match
	for _, v := range students {
		matches = append(matches, &Match{
			Student: v,
		})
	}

	var notFoundEmailReportCard []string
	for _, card := range cards {
		matched := false
		cardName := strings.ReplaceAll(card.Name, " ", "")
		for _, match := range matches {
			if strings.Contains(cardName, match.Student.Name) {
				if match.ReportCard != nil {
					return nil, fmt.Errorf("duplicate report card name: %s and %s", card.Name, match.ReportCard.Name)
				}
				match.ReportCard = card
				matched = true
			}
		}
		if !matched {
			notFoundEmailReportCard = append(notFoundEmailReportCard, card.Name)
		}
	}

	if len(notFoundEmailReportCard) > 0 {
		//return nil, fmt.Errorf("report card %s not found student", strings.Join(notFoundEmailReportCard, ","))
		fmt.Printf("report card %s not found student\n", strings.Join(notFoundEmailReportCard, ","))
	}

	var notFoundReportCardStudentNames []string
	for _, v := range matches {
		if v.ReportCard == nil {
			notFoundReportCardStudentNames = append(notFoundReportCardStudentNames, v.Student.Name)
		}
	}
	if len(notFoundReportCardStudentNames) > 0 {
		return nil, fmt.Errorf("student %s not found report card", strings.Join(notFoundReportCardStudentNames, ","))
	}
	if len(students) != len(cards) {
		//return nil, fmt.Errorf("students number %d != cards number %d", len(students), len(cards))
		fmt.Printf("students number %d != cards number %d\n", len(students), len(cards))
	}
	return matches, nil
}

func readStudents(emailPath string) ([]*Student, error) {
	excel, err := excelize.OpenFile(emailPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", file, err)
	}
	sheet := excel.GetSheetList()[0]
	rows, err := excel.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	var students []*Student
	var isNotEmailNames []string
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if !util.IsEmail(row[2]) {
			isNotEmailNames = append(isNotEmailNames, fmt.Sprintf("%s-%s", row[1], row[2]))
		}
		students = append(students, &Student{
			Name:  row[1],
			Email: row[2],
		})
	}
	if len(isNotEmailNames) > 0 {
		return nil, fmt.Errorf("%s email format error", strings.Join(isNotEmailNames, ","))
	}
	return students, nil
}

type ReportCard struct {
	Name string
	Path string
}

func readReportCards(reportCardFolder string) ([]*ReportCard, error) {
	entities, err := os.ReadDir(reportCardFolder)
	if err != nil {
		return nil, err
	}
	var files []*ReportCard
	for _, entity := range entities {
		fi, err := entity.Info()
		if err != nil {
			return nil, err
		}

		files = append(files, &ReportCard{
			Name: entity.Name(),
			Path: filepath.Join(reportCardFolder, fi.Name()),
		})
	}
	return files, nil
}
