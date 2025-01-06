package filesystem

import (
	"fmt"
	"time"
)

type FileSystem struct {
	rootDir *directoryDescriptor
}

func (fs * FileSystem) Mkfs () {
	fmt.Println("Creation of file system...")
	fs.rootDir = &directoryDescriptor{}
	fs.rootDir.Init(0)
	fs.rootDir.Data = make(map[string]Descriptor)
	fs.rootDir.Data["."] = fs.rootDir
	fs.rootDir.Data[".."] = fs.rootDir
	fmt.Println("File system created")
}

func (fs* FileSystem) Create (fileName string) {
	id := int(time.Now().UnixNano())
	descriptor := &fileDescriptor{}
	descriptor.Init(id)

	fs.rootDir.Data[fileName] = descriptor
	descriptor.Data = make(map[int]*Block)
	fmt.Println("Create file:", fileName,"| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Ls() {
	fmt.Println("Hard links of currect directory:")
	for f, d := range fs.rootDir.Data {
		if d, isFile := d.(*fileDescriptor); isFile {
			fmt.Println("Name:",f ,"\t id:",d.Id)
		}
	}
}

func (fs *FileSystem) Stat(fileName string) {
	descriptor := fs.rootDir.Data[fileName].(*fileDescriptor)
	fmt.Println("Type:", descriptor.FileType,
				"\tId:",descriptor.Id,
				"\tHard links count:", descriptor.Nlink,
				"\tSize:", descriptor.Size,
				"\tBlocks:", descriptor.Nblock)

}

func (fs *FileSystem) Link(linkWith, toLink string) {
	descriptor := fs.rootDir.Data[linkWith].(*fileDescriptor)
	descriptor.Nlink++
	fs.rootDir.Data[toLink] = descriptor
	fmt.Println("Create hard link", toLink, "with", linkWith)
}

func (fs *FileSystem) Unlink(fileName string) {
	fmt.Println("Delete file:", fileName)
	descriptor := fs.rootDir.Data[fileName].(*fileDescriptor)
	descriptor.Nlink--
	delete(fs.rootDir.Data, fileName)
}

func (fs *FileSystem) Find(fileName string) bool {
	return fs.rootDir.Data[fileName] != nil
}

func (fs *FileSystem) GetDescriptor(fileName string) *fileDescriptor {
	return fs.rootDir.Data[fileName].(*fileDescriptor)
}

func (fs *FileSystem) NullifyDescriptor(fileName string) {
    fs.rootDir.Data[fileName] = nil
}
