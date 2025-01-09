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
	descriptor := &FileDescriptor{}
	descriptor.Init(id)

	fs.RootDir.Data[fileName] = descriptor
	fmt.Println("Create file:", fileName, "| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Ls(dir *DirectoryDescriptor) {
	fmt.Println("Hard links of currect directory:")
	for f, d := range dir.Data {
		switch desc := d.(type) {
		case *FileDescriptor:
			fmt.Println("Name:", f, "\t id:", desc.Id, "\t type:", desc.FileType)
		case *DirectoryDescriptor:
			fmt.Println("Name:", f, "\t id:", desc.Id, "\t type:", desc.FileType)
		case *symlinkDescriptor:
			fmt.Println("Name:", f, "\t id:", desc.Id, "\t type:", desc.FileType)
		}
	}
}

func (fs *FileSystem) Stat(fileName string) {
	switch descriptor := fs.RootDir.Data[fileName].(type) {
	case *FileDescriptor:
		fmt.Println("Type:", descriptor.FileType,
			"\tId:", descriptor.Id,
			"\tHard links count:", descriptor.Nlink,
			"\tSize:", descriptor.Size,
			"\tBlocks:", descriptor.Nblock)
	case *DirectoryDescriptor:
		fmt.Println("Type:", descriptor.FileType,
			"\tId:", descriptor.Id,
			"\tHard links count:", descriptor.Nlink,
			"\tSize:", descriptor.Size,
			"\tBlocks:", descriptor.Nblock)
	case *symlinkDescriptor:
		fmt.Println("Type:", descriptor.FileType,
			"\tId:", descriptor.Id,
			"\tHard links count:", descriptor.Nlink,
			"\tSize:", descriptor.Size,
			"\tBlocks:", descriptor.Nblock)
	}
}

func (fs *FileSystem) Link(linkWith, toLink string) {
	descriptor := fs.RootDir.Data[linkWith].(*FileDescriptor)
	descriptor.Nlink++
	fs.RootDir.Data[toLink] = descriptor
	fmt.Println("Create hard link", toLink, "with", linkWith)
}

func (fs *FileSystem) Unlink(fileName string) {
	fmt.Println("Delete file:", fileName)
	descriptor := fs.RootDir.Data[fileName].(*FileDescriptor)
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

func (fs *FileSystem) Mkdir(prevDir *DirectoryDescriptor, dirName string) {
	id := int(time.Now().UnixNano())
	descriptor := &DirectoryDescriptor{}
	descriptor.Init(id)
	descriptor.Data["."] = descriptor
	if prevDir == fs.RootDir {
		descriptor.Data[".."] = fs.RootDir
	} else {
		descriptor.Data[".."] = prevDir
	}
	prevDir.Data[dirName] = descriptor
 	fmt.Println("Create directory:", dirName, "| Descriptor id:", descriptor.Id)
}

func (fs *FileSystem) Find(directory *DirectoryDescriptor, fileName string) bool {
	return directory.Data[fileName] != nil
}

func (fs *FileSystem) GetDescriptor(directory *DirectoryDescriptor, fileName string) Descriptor {
	return directory.Data[fileName]
}

func (fs *FileSystem) NullifyDescriptor(directory *DirectoryDescriptor, fileName string) {
	directory.Data[fileName] = nil
}
