package main

import (
	"flag"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
)

func removeAllFiles(ftp *sftp.Client, dir string) { //{{{
	files, err := ftp.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		isdir := file.IsDir()
		if !isdir {
			err := ftp.Remove(dir + file.Name())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(dir, file.Name(), " removed.")
			}
		}
	}
}                                         //}}}
func list(ftp *sftp.Client, dir string) { //{{{
	files, err := ftp.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Println(file.Mode(), file.ModTime(), file.Name(), file.Size())
	}
}                                              //}}}
func download(ftp *sftp.Client, file string) { //{{{
	fmt.Println("download file:" + file)
	srcf, err := ftp.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcf.Close()
	dstf, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
	}
	defer dstf.Close()

	if _, err = srcf.WriteTo(dstf); err != nil {
		fmt.Println(err)
	}
	fmt.Println("download completed")
}                                            //}}}
func upload(ftp *sftp.Client, file string) { //{{{
	fmt.Println("upload file:" + file)
	finfo, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := finfo.Size()
	srcf, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcf.Close()

	_, dfile := path.Split(file)

	dstf, err := ftp.Create(dfile)
	if err != nil {
		fmt.Println(err)
	}
	defer dstf.Close()

	for {
		var bufsize int64
		bufsize = 1024
		if filesize > 1024 {
			filesize -= 1024
		} else {
			bufsize = filesize
		}
		fmt.Print(".")
		buf := make([]byte, bufsize)
		n, _ := srcf.Read(buf)
		if n == 0 {
			break
		}
		dstf.Write(buf)
	}

	fmt.Println("upload completed")
}             //}}}
func main() { //{{{
	var ip, port, id, pwd, ufile, dfile, dir, drm string
	flag.StringVar(&ip, "h", "", "ip")
	flag.StringVar(&port, "port", "22", "port")
	flag.StringVar(&id, "u", "", "user")
	flag.StringVar(&pwd, "p", "", "password")
	flag.StringVar(&dfile, "d", "", "download file")
	flag.StringVar(&ufile, "up", "", "upload file")
	flag.StringVar(&dir, "l", "", "ls")
	flag.StringVar(&drm, "drm", "", "remove all files")
	flag.Parse()

	if ip == "" || id == "" || pwd == "" {
		flag.PrintDefaults()
		return
	}

	addr := ip + ":" + port
	config := &ssh.ClientConfig{
		User:            id,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
	}
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	ftp, err := sftp.NewClient(conn)
	if err != nil {
		panic("Failed to create client: " + err.Error())
	}
	// Close connection
	defer ftp.Close()

	if drm != "" {
		removeAllFiles(ftp, drm)
	}
	if dir != "" {
		list(ftp, dir)
	}
	if ufile != "" {
		upload(ftp, ufile)
	}
	if dfile != "" {
		download(ftp, dfile)
	}

} //}}}
