package filesystem

type OpenFileDescriptor struct {
	Id int
	Desc *fileDescriptor
	Offset int
}