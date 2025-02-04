package core

import (
	fs "OS_lab5/FileSystem"
	"fmt"
	"strings"
)

type Core struct {
	fs *fs.FileSystem
	openFileDescriptors []*fs.OpenFileDescriptor
	blockSize int
	Cwd *fs.DirectoryDescriptor
}

func (c *Core) Mkfs (descriptorsCount int) {
	fmt.Println("System initialization...")
	c.openFileDescriptors = make([]*fs.OpenFileDescriptor, descriptorsCount)
	c.blockSize = fs.BlockSize
	fmt.Println("Create core with", descriptorsCount, "possible open file descpriptors")
	c.fs = &fs.FileSystem{}
	c.fs.Mkfs()
	c.Cwd = c.fs.RootDir
	fmt.Println("System is ready to work!")
}

func (c *Core) Create(filePath string) {
	curDir, desc, fileName := c.lookup(filePath,false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", filePath)
		return
	}
	if (desc != nil) {
		fmt.Println("Error: File", filePath, "to create already exist")
		return
	}
	c.fs.Create(curDir, fileName)
}

func (c *Core) Ls(path ...string) {
	if path == nil {
		c.fs.Ls(c.Cwd)
	} else {
		_, dir, _ := c.lookup(path[0], false)
		if (dir == nil) {
			fmt.Println("Error: Directory", path[0], "does not exist")
			return
		}
		fmt.Println("List of hard links", path[0], "directory")
		c.fs.Ls(dir.(*fs.DirectoryDescriptor))
	}
}

func (c *Core) Stat(filePath string) {
	curDir, desc, _ := c.lookup(filePath, false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", filePath)
		return
	}
	if (desc == nil) {
		fmt.Println("Error: File", filePath," to stat does not exist")
		return
	}
	fmt.Println("Print stat for file:", filePath)
	c.fs.Stat(desc)
}

func (c *Core) Link(linkWithPath, toLinkPath string) {
	curDirLinkWith, descLinkWith, fileLinkWith := c.lookup(linkWithPath, false)
	if (curDirLinkWith == nil) {
		fmt.Println("Error: incorrect path", linkWithPath)
		return
	}
	if (descLinkWith == nil) {
		fmt.Println("Error: File", linkWithPath, "link with does not exist exist")
		return
	}

	curDirToLink, descToLink, fileToLink := c.lookup(toLinkPath, false)
	if (curDirToLink == nil) {
		fmt.Println("Error: incorrect path", toLinkPath)
		return
	}
	if (descToLink != nil) {
		fmt.Println("Error: File", toLinkPath, "to link already exist")
		return
	}
	desc := descLinkWith.(*fs.FileDescriptor)
	c.fs.Link(curDirToLink, desc, fileToLink)
	fmt.Println("Create hard link", fileToLink, "with", fileLinkWith)
}

func (c *Core) Unlink(filePath string) {
	curDir, desc, fileName := c.lookup(filePath, false)
	_, isDir := desc.(*fs.DirectoryDescriptor)
	if (isDir) {
		fmt.Println("Error: can not unlink hard link to directory", filePath)
		return
	}
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", filePath)
		return
	}
	if (desc == nil) {
		fmt.Println("Error: File", filePath, "to delete does not exist exist")
		return
	}
	fileDesc := desc.(*fs.FileDescriptor)
	c.fs.Unlink(curDir, fileName)
	fmt.Println("Delete file:", filePath)
	if (fileDesc.Nlink == 0 && fileDesc.NOpen == 0) {
		c.fs.NullifyDescriptor(curDir, fileName)
	}
}

func (c *Core) Open(filePath string) *fs.OpenFileDescriptor{
	curDir, desc, _ := c.lookup(filePath, true)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", filePath)
		return nil
	}
	if (desc == nil) {
		fmt.Println("Error: File", filePath, "to open does not exist")
		return nil
	}
	index := c.findFreeIndex()
	if (index == -1) {
		fmt.Println("No free descriptor indexes")
		return nil
	}

	fmt.Println("Open file", filePath)
	fileDesc := desc.(*fs.FileDescriptor)
	fileDesc.NOpen++
	openFileDescriptor := &fs.OpenFileDescriptor{Desc: fileDesc, Offset: 0, Id: index}
	c.openFileDescriptors[index] = openFileDescriptor
	openFileDescriptor.Desc.Data = make(map[int]*fs.Block)
	return openFileDescriptor
}

func (c *Core) findFreeIndex() int {
	freeIndex := -1
	for i, v := range c.openFileDescriptors {
		if (v == nil) {
			freeIndex = i
			break
		}
	}
	return freeIndex
}


func (c *Core) Close(fd *fs.OpenFileDescriptor) *fs.OpenFileDescriptor {
	if (fd == nil) {
		fmt.Println("Error: closing of non-existing file")
		return nil
	}
	fmt.Println("Closing file")
	c.openFileDescriptors[fd.Id] = nil
	fd.Desc.NOpen--
	if(fd.Desc.Nlink == 0 && fd.Desc.NOpen == 0) {
		fd.Desc = nil
	}
	return nil
}

func (c *Core) Truncate(filePath string, size int) {
	if (size <= 0) {
		fmt.Println("Error: Incorrect size to truncate, must be bigger than 0")
		return
	}
	curDir, desc, _ := c.lookup(filePath, true)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", filePath)
		return
	}
	if (desc == nil) {
		fmt.Println("Error: File", filePath, "does not exist")
		return 
	}
	fileDesc := desc.(*fs.FileDescriptor)
	if (fileDesc.Size > size) {
		newBlockCount := size / c.blockSize
		remainingBytes := size % c.blockSize
		if (remainingBytes > 0) {
			newBlockCount++
		}
		for i := newBlockCount; fileDesc.Nblock > newBlockCount; i++ {
			if (fileDesc.Data[i] == nil) {
				continue
			}
			delete(fileDesc.Data, i)
			fileDesc.Nblock--
		}
	}
	fileDesc.Size = size
}

func (c *Core) Read(fd *fs.OpenFileDescriptor, size int) {
	if (size <= 0) {
		fmt.Println("Error: Incorrect size to read, must be bigger than 0")
		return
	}
	if (size > fd.Desc.Size) {
		fmt.Println("Error: Incorrect size to read, must not be bigger than file size")
		return
	}
	curOffset := fd.Offset
	totalSize := size
	bytesToRead := 0
	res := ""
	for totalSize > 0 {
		curBlock := curOffset / c.blockSize
		offsetInsideBlock := curOffset % c.blockSize
		if (totalSize > (c.blockSize - offsetInsideBlock)) {
			bytesToRead = c.blockSize - offsetInsideBlock
		} else {
			bytesToRead = totalSize
		}
		if (fd.Desc.Data[curBlock] == nil) {
			for i := 0; i < bytesToRead; i++ {
				res += "0"
			}
			curOffset += bytesToRead
			totalSize -= bytesToRead
			continue
		}
		block := fd.Desc.Data[curBlock]
		readTo := offsetInsideBlock + bytesToRead
		res += string(block[offsetInsideBlock:readTo])
		curOffset += bytesToRead
		totalSize -= bytesToRead
	}
	fmt.Println(res)
}

func (c *Core) Write(fd *fs.OpenFileDescriptor, data []byte) {
	totalSize := len(data)
	if (totalSize > fd.Desc.Size) {
		fmt.Println("Error: Incorrect size to write, must be less than file size")
		return
	}
	curOffset := fd.Offset
	bytesToWrite := 0
	for totalSize > 0 {
		curBlock := curOffset / c.blockSize
		offsetInsideBlock := curOffset % c.blockSize
		if (fd.Desc.Data[curBlock] == nil) {
			block := new(fs.Block)
			fd.Desc.Data[curBlock] = block
			fd.Desc.Nblock = len(fd.Desc.Data)
		}
		if (totalSize > (c.blockSize - offsetInsideBlock)) {
			bytesToWrite = c.blockSize - offsetInsideBlock
		} else {
			bytesToWrite = totalSize
		}
		block := fd.Desc.Data[curBlock]
		writeTo := offsetInsideBlock + bytesToWrite
		getDataFrom := curOffset - fd.Offset
		getDataTo := getDataFrom + bytesToWrite
		copy(block[offsetInsideBlock:writeTo], data[getDataFrom:getDataTo])
		curOffset += bytesToWrite
		totalSize -= bytesToWrite
	}
}

func (c *Core) Seek(fd *fs.OpenFileDescriptor, offset int) {
	if (offset < 0) {
		fmt.Println("Error: Offset can not be less than 0")
		return
	}
	if (offset > fd.Desc.Size) {
		fmt.Println("Error: Offset can not be bigger tnah file size")
		return
	}
	fd.Offset = offset
}

func (c *Core) Symlink(linkPath, content string) {
	if (len(content) > 32) {
		fmt.Println("Error: symlink content can not be bigger than block size", c.blockSize)
		return
	}
	curDir, linkDesc, linkName := c.lookup(linkPath, false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", linkPath)
		return
	}
	if (linkDesc != nil) {
		fmt.Println("Error: Directory", linkPath, "already exist")
		return
	}
	c.fs.Symlink(curDir, linkName, content)
	fmt.Println("Create symlink:", linkPath, " to file", content)
}

func (c *Core) Mkdir(path string) {
	curDir, dir, dirName := c.lookup(path, false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", path)
		return
	}
	if (dir != nil) {
		fmt.Println("Error: Directory", path, "already exist")
		return
	}
	c.fs.Mkdir(curDir, dirName)
	fmt.Println("Create directory:", path)
}

func (c *Core) Rmdir(path string) {
	if (path == "/") {
		fmt.Println("Cannot delete root directory, don't play with rm -rf /")
		return
	}
	curDir, dir, dirName := c.lookup(path, false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", path)
		return
	}
	if (dir == nil) {
		fmt.Println("Error: Directory to delete", path, "does not exist")
		return
	}
	dirContent := dir.(*fs.DirectoryDescriptor).Data
	if (len(dirContent) > 2) {
		fmt.Println("Error: Directory to delete", path,"is not empty")
		return
	}
	c.fs.Rmdir(curDir, dirName)
	fmt.Println("Delete directory", path)
}

func (c *Core) Cd(path string) {
	curDir, dir, _ := c.lookup(path, false)
	if (curDir == nil) {
		fmt.Println("Error: incorrect path", path)
		return
	}
	if (dir == nil) {
		fmt.Println("Error: Directory to change to current", path, "does not exist")
		return
	}
	c.Cwd = dir.(*fs.DirectoryDescriptor)
	fmt.Println("Change current working directory to", path)
}

func (c *Core) lookup(pathname string, follow bool) (*fs.DirectoryDescriptor, fs.Descriptor, string) {
	curDir := c.Cwd
	pathComponents := strings.Split(pathname, "/")
	if pathname[0] == '/' {
		curDir = c.fs.RootDir
		pathComponents = pathComponents[1:]
	}
	lastComponentIdx := len(pathComponents) - 1
	for i := 0; i < len(pathComponents); i++ {
		componentName := pathComponents[i]
		if !c.fs.Find(curDir, componentName){
			if i == lastComponentIdx {
				return curDir, nil, componentName
			}
			return nil, nil, ""
		}
		desc := c.fs.GetDescriptor(curDir, componentName)
		if symlink, isSymlink := desc.(*fs.SymlinkDescriptor); isSymlink {
			if i == lastComponentIdx && !follow {
				fmt.Println("Error: last component of pathname cannot be symlink")
				return nil, nil, ""
			}
			linkedPath := strings.Split(symlink.Data, "/")
			pathComponents = append(pathComponents[:i], append(linkedPath, pathComponents[i+1:]...)...)
			lastComponentIdx = len(pathComponents) - 1
			i--
			continue
		}
		if i != lastComponentIdx {
			if nextDir, isDir := desc.(*fs.DirectoryDescriptor); isDir {
				curDir = nextDir
			} else {
				return curDir, nil, ""
			}
		} else {
			return curDir, desc, componentName
		}
	}
	return nil, nil, ""
}
