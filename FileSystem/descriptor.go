package filesystem

type Descriptor interface {
	Init(fileType string, id int)
}

type fileDescriptor struct{
	FileType string
	Nlink, NOpen, Size, Id int
	Data map[int]*Block
	Nblock int
}

func (fd *fileDescriptor) Init (fileType string, id int) {
	fd.FileType = "reg"
	fd.Id = id
	fd.Nlink = 1
}
