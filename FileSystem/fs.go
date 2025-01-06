package filesystem

import (
	"fmt"
	"time"
)

type FileSystem struct {
	RootDir *DirectoryDescriptor
}

func (fs *FileSystem) Mkfs() {
	fmt.Println("Creation of file system...")
	fs.RootDir = &DirectoryDescriptor{}
	fs.RootDir.Init(0)
	fs.RootDir.Data["."] = fs.RootDir
	fs.RootDir.Data[".."] = fs.RootDir
	fmt.Println("File system created")
}

func (fs *FileSystem) Create(fileName string) {
	id := int(time.Now().UnixNano())
	descriptor := &fileDescriptor{}
	descriptor.Init(id)

	fs.RootDir.Data[fileName] = descriptor
	fmt.Println("Create file:", fileName, "| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Ls() {
	fmt.Println("Hard links of currect directory:")
	for f, d := range fs.RootDir.Data {
		if d, isFile := d.(*fileDescriptor); isFile {
			fmt.Println("Name:", f, "\t id:", d.Id)
		}
	}
}

func (fs *FileSystem) Stat(fileName string) {
	descriptor := fs.RootDir.Data[fileName].(*fileDescriptor)
	fmt.Println("Type:", descriptor.FileType,
		"\tId:", descriptor.Id,
		"\tHard links count:", descriptor.Nlink,
		"\tSize:", descriptor.Size,
		"\tBlocks:", descriptor.Nblock)

}

func (fs *FileSystem) Link(linkWith, toLink string) {
	descriptor := fs.RootDir.Data[linkWith].(*fileDescriptor)
	descriptor.Nlink++
	fs.RootDir.Data[toLink] = descriptor
	fmt.Println("Create hard link", toLink, "with", linkWith)
}

func (fs *FileSystem) Unlink(fileName string) {
	fmt.Println("Delete file:", fileName)
	descriptor := fs.RootDir.Data[fileName].(*fileDescriptor)
	descriptor.Nlink--
	delete(fs.RootDir.Data, fileName)
}

func (fs *FileSystem) Symlink(linkname, content string) {
	id := int(time.Now().UnixNano())
	descriptor := &symlinkDescriptor{}
	descriptor.Init(id)
	descriptor.Data = content
	fs.RootDir.Data[linkname] = descriptor
	fmt.Println("Create symlink:", linkname, " to file", content, "| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Find(fileName string) bool {
	return fs.RootDir.Data[fileName] != nil
}

func (fs *FileSystem) GetDescriptor(fileName string) *fileDescriptor {
	return fs.RootDir.Data[fileName].(*fileDescriptor)
}

func (fs *FileSystem) NullifyDescriptor(fileName string) {
	fs.RootDir.Data[fileName] = nil
}
