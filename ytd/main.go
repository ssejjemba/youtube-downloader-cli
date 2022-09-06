package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	. "youtube_downloader/provider"
)

const usageString = `Usage: ytd [OPTION] [URL]
Download a video from youtube.
Example: ytd -o "Chasing The Sky".mp4 https://youtu.be/Tds5IfWpXo0`

func main() {
	flag.Usage = func ()  {
		fmt.Println(usageString)
		flag.PrintDefaults()
	}
	usr, _ := user.Current()
	var outputFile string
	flag.StringVar(&outputFile, "o", "sample.mp4", "The output file name")
	var outputDir string
	flag.StringVar(&outputDir, "d", filepath.Join(usr.HomeDir, "Movies", "ytd"), "The output directory")
	flag.Parse()
	log.Println(flag.Args())

	log.Println("download to dir=", outputDir)
	y := NewYoutube(true)
	arg := flag.Arg(0)

	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("error found: ", err)
	}
	
	if err := y.StartDownload(filepath.Join(outputDir, outputFile)); err != nil {
		fmt.Println("err:", err)
	}

}