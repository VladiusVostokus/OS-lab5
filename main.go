package main

import (
	c "OS_lab5/Core"
	"fmt"
)

func main() {
	core := c.Core{}
	core.Mkfs(100)

	fmt.Println("\n=====================Test creation of FS and files=========================")
	core.Create("file.txt")
	core.Create("file.txt")
	core.Create("a.txt")
	core.Ls()
	core.Stat("file.txt")
	core.Stat("bbbbb.txt")

	fmt.Println("\n===========================Test link/unlink================================")
	core.Link("file.txt","file.txt")
	core.Link("file3123.txt","file1.txt")
	core.Link("file.txt","file2.txt")
	core.Stat("file.txt")
	core.Stat("file2.txt")
	core.Unlink("fileaaaa.txt")
	core.Unlink("file.txt")
	core.Stat("file2.txt")

	fmt.Println("\n============================Test open/close================================")
	fd := core.Open("file2.txt")
	fd = core.Close(fd)

	fmt.Println("\n============================Test truncate==================================")
	core.Truncate("file2.txt",-10)
	core.Truncate("file2.txt",10)
	core.Stat("file2.txt")

	fmt.Println("\n==============================Test write/read==============================")
	fd = core.Open("file2.txt")
	str1 := []byte("20 len str is here !")
	str2 := []byte("10 len str")
	core.Write(fd, str1)
	core.Write(fd, str2)
	core.Read(fd, 10)
	core.Read(fd, 20)

	fmt.Println("\n=======================Test write/read with offset=========================")
	core.Truncate("file2.txt",40)
	str := []byte("This string contains 32 symbols 35!")
	core.Write(fd, str)
	core.Read(fd, 32)
	core.Read(fd, 35)
	core.Seek(fd, 40)
	core.Read(fd, 32)
	core.Seek(fd, 5)
	core.Read(fd, 5)
	core.Seek(fd, 30)
	core.Read(fd, 5)

	fmt.Println("====================Test offset < 0 and offset < size======================")
	core.Seek(fd, -10)
	core.Seek(fd, 100)
	fd = core.Close(fd)

	fmt.Println("\n=======================Test truncate with less size========================")
	fd = core.Open("file2.txt")
	core.Truncate("file2.txt", 65)
	str = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	core.Write(fd, str)
	core.Stat("file2.txt")
	core.Read(fd, 55)
	core.Truncate("file2.txt", 10)
	core.Read(fd, 10)
	core.Stat("file2.txt")
	core.Truncate("file2.txt", 200)
	core.Seek(fd, 100)
	core.Read(fd, 4)
	fd = core.Close(fd)

	fmt.Println("\n===================Test write/read after open and unlink===================")
	core.Create("unlink.txt")
	core.Truncate("unlink.txt", 23)
	fdd := core.Open("unlink.txt")
	fdd2 := core.Open("unlink.txt")
	core.Unlink("unlink.txt")
	core.Stat("unlink.txt")
	str = []byte("Content of deleted file")
	core.Write(fdd, str)
	core.Read(fdd, 23)
	fdd = core.Close(fdd)
	fdd = core.Open("unlink.txt")

	core.Seek(fdd2, 3)
	aaa := []byte("aaa")
	core.Write(fdd2, aaa)
	core.Read(fdd2, 20)
	fdd2 = core.Close(fdd2)

	fmt.Println("\n==============================Test symlinks================================")
	core.Symlink("symlink0.txt","contentBiggerThanBlockSize.txt")
	core.Symlink("symlink.txt","somefile.txt")
	core.Stat("symlink.txt")

	fmt.Println("\n==============================Test directories=============================")
	core.Mkdir("/a")
	core.Mkdir("/a/b")
	core.Mkdir("/a/b/c")
	core.Mkdir("/a/b/c")
	core.Mkdir("/a/b/c/d/e")
	core.Ls()
	core.Cd("/a/b")
	core.Ls("/a/b/c")
	core.Rmdir("/a/b/c")
	core.Ls("/a/b/c")
	core.Ls("/a/././b/../b")
	core.Mkdir("c")
	core.Create("c/f.txt")
	core.Truncate("c/f.txt", 10)
	fd3 := core.Open("c/f.txt")
	aaa2 := []byte("aaa2")
	core.Write(fd3, aaa2)
	core.Read(fd3, 4)
	core.Link("c/f.txt","c/f2.txt")
	core.Unlink("c/f.txt")
	core.Ls("c")
	core.Unlink("/a")
	core.Unlink("/a/.")
	core.Unlink("/a/..")
	core.Stat("/a")
	core.Create("/a/bbb.txt")
	core.Cd("/a/bbb.txt/c")

	fmt.Println("\n======================Test symlinks to dirs and files=======================")
	core.Symlink("/a/symlink.txt","b/c")
	core.Mkdir("/a/b/c/d")
	core.Mkdir("/a/b/c/d/e")
	core.Cd("/a/symlink.txt/d")
	core.Symlink("/a/b/c/d/eee","e")
	core.Cd("/a/symlink.txt/d/eee")
	core.Create("fileInD.txt")
	core.Truncate("fileInD.txt",10)
	core.Symlink("/a/b/fileDSym","c/d/fileInD.txt")
	fd4 := core.Open("/a/b/fileDSym")
	core.Write(fd4,[]byte("aaa"))
}