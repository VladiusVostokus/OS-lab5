package filesystem

type Descriptor interface {
	Init(id int)
}

type fileDescriptor struct{
	FileType string
	Nlink, NOpen, Size, Id int
	Data map[int]*Block
	Nblock int
}

func (fd *fileDescriptor) Init (id int) {
	fd.FileType = "reg"
	fd.Id = id
	fd.Nlink = 1
	fd.Data = make(map[int]*Block)
}

type symlinkDescriptor struct {
	FileType string
	Nlink, NOpen, Size, Id int
	Data string
	Nblock int
}

func (fd *symlinkDescriptor) Init (id int) {
	fd.FileType = "sym"
	fd.Id = id
	fd.Nlink = 1
}

type directoryDescriptor struct {
	FileType string
	Nlink, NOpen, Size, Id int
	Data map[string]Descriptor
	Nblock int
}

func (fd *directoryDescriptor) Init (id int) {
	fd.FileType = "dir"
	fd.Id = id
	fd.Nlink = 1
	fd.Data = make(map[string]Descriptor)
}
