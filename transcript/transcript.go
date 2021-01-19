package transcript

//成绩单
type Transcript struct {
	Name           string
	Class          string
	Grades         []string
	StudentComment string
	ParentComment  string
	TeacherComment string
}

type Transcripts []Transcript
