package filesystem

type fileDescriptor struct{
	FileType string
	Nlink, NOpen, Size, Id int
	Data map[int]*Block
	Nblock int
}