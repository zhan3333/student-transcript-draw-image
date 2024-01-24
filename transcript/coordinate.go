package transcript

// 坐标配置

type Coordinate struct {
	X int
	Y int
}

// 学生姓名坐标
var studentNameCoordinate = Coordinate{510, 390}

//var studentNameCoordinate = Coordinate{1000, 920}

// 班级坐标
var classCoordinate = Coordinate{980, 390}

// 童言妙语
var studentCommentCoordinate = Coordinate{260, 1305}

// 教师评语
var teacherCommentCoordinate = Coordinate{260, 940}

// 家长寄语
var parentCommentCoordinate = Coordinate{880, 1305}

// 成绩起始点与偏移量
var gradesCoordinate = Coordinate{470, 1750}

// 成绩之间的偏移量
var gradesOffset = 142
