package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	. "youtube_downloader/provider"
)

func main() {
	flag.Parse()
	log.Println(flag.Args())
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/ytd", usr.HomeDir)

	log.Println("download to dir=", currentDir)
	y := NewYoutube(true)
	arg := flag.Arg(0)

	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("error found: ", err)
	}
	
	if err := y.StartDownload(filepath.Join(currentDir, "sample.mp4")); err != nil {
		fmt.Println("err:", err)
	}

}