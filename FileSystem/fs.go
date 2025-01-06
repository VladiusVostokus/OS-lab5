package filesystem

import (
	"fmt"
	"time"
)

type FileSystem struct {
	directory map[string]Descriptor
}

func (fs * FileSystem) Mkfs () {
	fmt.Println("Creation of file system...")
	fs.directory = make(map[string]Descriptor)
	fmt.Println("File system created")
}

func (fs* FileSystem) Create (fileName string) {
	id := int(time.Now().UnixNano())
	descriptor := &fileDescriptor{}
	descriptor.Init("reg", id)

	fs.directory[fileName] = descriptor
	descriptor.Data = make(map[int]*Block)
	fmt.Println("Create file:", fileName,"| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Ls() {
	fmt.Println("Hard links of currect directory:")
	for f, d := range fs.directory {
		if d, isFile := d.(*fileDescriptor); isFile {
			fmt.Println("Name:",f ,"\t id:",d.Id)
		}
	}
}

func (fs *FileSystem) Stat(fileName string) {
	descriptor := fs.directory[fileName].(*fileDescriptor)
	fmt.Println("Type:", descriptor.FileType, 
				"\tId:",descriptor.Id, 
				"\tHard links count:", descriptor.Nlink, 
				"\tSize:", descriptor.Size,
				"\tBlocks:", descriptor.Nblock)
	
}

func (fs *FileSystem) Link(linkWith, toLink string) {
	descriptor := fs.directory[linkWith].(*fileDescriptor)
	descriptor.Nlink++
	fs.directory[toLink] = descriptor
	fmt.Println("Create hard link", toLink, "with", linkWith)
}

func (fs *FileSystem) Unlink(fileName string) {
	fmt.Println("Delete file:", fileName)
	descriptor := fs.directory[fileName].(*fileDescriptor)
	descriptor.Nlink--
	delete(fs.directory, fileName)
}

func (fs *FileSystem) Find(fileName string) bool {
	return fs.directory[fileName] != nil
}

func (fs *FileSystem) GetDescriptor(fileName string) *fileDescriptor {
	return fs.directory[fileName].(*fileDescriptor)
}

func (fs *FileSystem) NullifyDescriptor(fileName string) {
    fs.directory[fileName] = nil
}
