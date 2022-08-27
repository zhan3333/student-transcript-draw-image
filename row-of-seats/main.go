package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

var students []Student
var room [][]Seat

func init() {
	for i := 0; i < 7; i++ {
		room = append(room, make([]Seat, 8))
		for j := 0; j < len(room[i]); j++ {
			if i == 6 && (j == 0 || j == 1 || j == 6 || j == 7) {
				room[i][j].Enable = false
			} else {
				room[i][j].Enable = true
			}
			room[i][j].Row = i
			room[i][j].Line = j
		}
	}
}

func main() {
	if err := initStudent("row-of-seats/学生性别成绩表.xlsx"); err != nil {
		panic(err)
	}
	fmt.Printf("读取到了 %d 位同学\n", len(students))
	students = sortStudents(students)
	for _, s := range students {
		fmt.Printf("%s\n", s)
	}
	seatArrangement()
	if err := printSeatsToFile("row-of-seats/座位表.xlsx"); err != nil {
		fmt.Printf("输出座位表失败: %w\n", err)
	}
}

func sortStudents(students []Student) []Student {
	for i := 0; i < len(students)-1; i++ {
		max := i
		for j := i + 1; j < len(students); j++ {
			if students[j].Scope > students[max].Scope {
				max = j
			}
		}
		if max != i {
			students[i], students[max] = students[max], students[i]
		}
	}
	return students
}

func initStudent(filename string) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("读取成绩文件失败: %w", err)
	}
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return fmt.Errorf("读取 excel 第一个 sheet 失败: %w", err)
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if row[0] == "" || row[2] == "" {
			return fmt.Errorf("读取第 %d 行失败: 有数据未填写", i+1)
		}
		scope, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return fmt.Errorf("读取第 %d 行总分失败: %w", i+1, err)
		}
		sex := None
		if row[1] == "女" {
			sex = Female
		}
		if row[1] == "男" {
			sex = Male
		}
		top10 := false
		level := ""
		if i >= 1 && i <= 10 {
			top10 = true
			level = "top"
		}
		if i > 10 && i <= 20 {
			level = "top"
		}
		if i > 20 && i <= 40 {
			level = "mid"
		}
		if i > 40 {
			level = "bottom"
		}
		students = append(students, Student{
			Name:  row[0],
			Scope: scope,
			Sex:   sex,
			Top10: top10,
			Level: level,
		})
	}
	return nil
}

// 安排座位
func seatArrangement() {
	//1. 前10名坐一起
	//2. 11-50名按照排名分为 上中下
	//3. 上中、中下 搭配
	//4. 男女搭配
	//5. 中下坐 1，2排
	//6. 上中坐 3, 4, 5, 6 排

	var count int

	for len(students) > 0 && count < 10 {
		var students2 []Student

		for _, student := range students {
			if !setToSeat(student) {
				// 没找到座位
				students2 = append(students2, student)
			}
		}
		students = students2
		count++
		if len(students) > 0 {
			fmt.Printf("%d 位同学第 %d 次没找到座位\n", len(students), count)
		}
	}
	if len(students) == 0 {
		fmt.Println("座位安排完毕")
		printSeats()
		return
	} else {
		fmt.Println("还有些同学直接安排到空座位上")
	}

	// 安排剩下的同学到空座位
	if len(students) > 0 {
		var students2 []Student
		for _, student := range students {
			if !setToEmpty(student) {
				students2 = append(students2, student)
			}
		}
		students = students2
		if len(students) > 0 {
			fmt.Printf("%d 位同学在安排空座位后还没有找到座位\n", len(students))
		}
	}
	if len(students) == 0 {
		fmt.Println("座位安排完毕")
		printSeats()
		return
	} else {
		fmt.Println("还有这些同学实在找不到位置了，手动填吧")
		for _, student := range students {
			fmt.Println(student)
		}
	}
}

// 打印目前安排的座位
func printSeats() {
	var data [][]string

	for i := range room {
		var row []string
		for j := range room[i] {
			seat := room[i][j]
			if seat.Student != nil {
				row = append(row, fmt.Sprintf("%s(%s%s)",
					seat.Student.Name,
					seat.Student.Sex.String(),
					seat.Student.Level,
				))
			} else {
				row = append(row, "empty")
			}
		}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"一", "二", "三", "四", "五", "六", "七", "八"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

// 输出到 excel 中
func printSeatsToFile(filename string) error {
	var sheet = "Sheet1"
	var err error
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet(sheet)
	// 设置单元格的值
	err = f.SetSheetRow(sheet, "A1", &[]string{
		"一", "二", "三", "四", "五", "六", "七", "八",
	})
	if err != nil {
		return err
	}

	for i := range room {
		var row []string
		for j := range room[i] {
			seat := room[i][j]
			if seat.Student != nil {
				row = append(row, seat.Student.Name)
			} else {
				row = append(row, "空")
			}
		}
		if err = f.SetSheetRow(sheet, fmt.Sprintf("A%d", i+2), &row); err != nil {
			return err
		}
	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return nil
}

// 按照规则给学生安排座位
func setToSeat(student Student) bool {
	if student.Top10 {
		// 从第三行第四行找位置
		for setRow := 2; setRow <= 3; setRow++ {
			for setLine := 0; setLine < 8; setLine++ {
				seat := room[setRow][setLine]
				if !seat.Enable || seat.Student != nil {
					continue
				}
				if setLine&1 == 0 {
					side := room[setRow][setLine+1]
					if side.Student == nil || side.Student.Sex != student.Sex {
						// 如果邻桌为空或者性别不一样，则可以坐
						room[setRow][setLine].Student = &student
						return true
					}
				} else {
					side := room[setRow][setLine-1]
					if side.Student == nil || side.Student.Sex != student.Sex {
						// 如果邻桌为空或者性别不一样，则可以坐
						room[setRow][setLine].Student = &student
						return true
					}
				}
			}
		}
	}
	if student.Level == "top" {
		// 上从 3-7 排找位置
		for setRow := 2; setRow < 7; setRow++ {
			for setLine := 0; setLine < 8; setLine++ {
				seat := room[setRow][setLine]
				if !seat.Enable || seat.Student != nil {
					continue
				}
				if setLine&1 == 0 {
					// 偶数行，看同桌是否为中等
					side := room[setRow][setLine+1]
					if side.Student == nil || (side.Student.Level == "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				} else {
					// 奇数行，看同桌是否为中等
					side := room[setRow][setLine-1]
					if side.Student == nil || (side.Student.Level == "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				}
			}
		}
	}
	if student.Level == "bottom" {
		// 下从第一行第二行找位置
		for setRow := 0; setRow <= 1; setRow++ {
			for setLine := 0; setLine < 8; setLine++ {
				seat := room[setRow][setLine]
				if !seat.Enable || seat.Student != nil {
					continue
				}
				if setLine&1 == 0 {
					// 偶数行，看同桌是否为中等
					side := room[setRow][setLine+1]
					if side.Student == nil || (side.Student.Level == "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				} else {
					// 奇数行，看同桌是否为中等
					side := room[setRow][setLine-1]
					if side.Student == nil || (side.Student.Level == "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				}
			}
		}
	}
	if student.Level == "mid" {
		// 中在 1-7 都能坐
		// 安排在同桌为空，或者同桌为上/下
		for setRow := 0; setRow < 7; setRow++ {
			for setLine := 0; setLine < 8; setLine++ {
				seat := room[setRow][setLine]
				if !seat.Enable || seat.Student != nil {
					continue
				}
				if setLine&1 == 0 {
					// 偶数行，看同桌是否为不为中等
					side := room[setRow][setLine+1]
					if side.Student == nil || (side.Student.Level != "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				} else {
					// 奇数行，看同桌是否不为中等
					side := room[setRow][setLine-1]
					if side.Student == nil || (side.Student.Level != "mid" && side.Student.Sex != student.Sex) {
						room[setRow][setLine].Student = &student
						return true
					}
				}
			}
		}
	}
	// 没找到座位，等下次安排
	return false
}

// 将学生安排到空座位上
func setToEmpty(student Student) bool {
	for i := range room {
		for j := range room[i] {
			seat := room[i][j]
			if !seat.Enable || seat.Student != nil {
				continue
			}
			var side Seat
			// 偶数右边为同桌，奇数左边为同桌
			if j&1 == 0 {
				side = room[i][j+1]
			} else {
				side = room[i][j-1]
			}
			if side.Student == nil ||
				side.Student.Sex == None ||
				(side.Student.Sex != student.Sex) {
				// 邻桌为空或者无性别或者性别不一样就可以安排坐
				room[i][j].Student = &student
				return true
			}
		}
	}
	return false
}

type Sex int

func (s Sex) String() string {
	switch s {
	case Male:
		return "男"
	case Female:
		return "女"
	default:
		return "无"
	}
}

var (
	None   Sex = 0
	Male   Sex = 1
	Female Sex = 2
)

type Student struct {
	Name  string
	Scope float64
	Sex   Sex
	// top, mid, bottom
	Level string
	Top10 bool
}

func (s Student) String() string {
	return fmt.Sprintf("%s-%f", s.Name, s.Scope)
}

type Seat struct {
	Student *Student
	// 列, 总共 8 列 0-7
	Line int
	// 行，从第 0 行开始，逐渐往后排 0-7
	Row int
	// 座位是否启用
	Enable bool
}
