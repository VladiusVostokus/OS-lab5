package filesystem

type fileDescriptor struct{
	FileType string
	Nlink, Size, Id int
	Data map[int]*Block
	Nblock int
	IsOpen bool
}