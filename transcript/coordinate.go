package transcript

// 坐标配置

type Coordinate struct {
	X int
	Y int
}

// 学生姓名坐标
var studentNameCoordinate = Coordinate{1000, 920}

// 班级坐标
var classCoordinate = Coordinate{2000, 920}

// 童言妙语
var studentCommentCoordinate = Coordinate{750, 2655}

// 教师评语
var teacherCommentCoordinate = Coordinate{730, 1875}

// 家长寄语
var parentCommentCoordinate = Coordinate{1930, 2570}

// 成绩起始点与偏移量
var gradesCoordinate = Coordinate{1000, 3750}

// 成绩之间的偏移量
var gradesOffset = 300
