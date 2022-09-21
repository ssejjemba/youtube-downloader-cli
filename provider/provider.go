package provider

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

func NewYoutube(debug bool) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercentage: make(chan int64, 100)}
}

type stream map[string]string


// TODO: Use consistent casing
type Youtube struct {
	DebugMode bool
	StreamList []stream
	VideoID string
	VideoInfo string
	DownloadPercentage chan int64
	contentLength float64
	totalWrittenBytes float64
	downloadLevel float64
}

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoId(url)
	if err != nil {
		return fmt.Errorf("findvideoID error=%s", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}

	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}

	return nil
}

func (y *Youtube) StartDownload(destFile string) error {
	var err error
	for _, v := range y.StreamList {
		url := v["url"] + "&signature=" + v["sig"]
		y.log(fmt.Sprintln("Download url=", url))

		y.log(fmt.Sprintln("Download to file=", destFile))
		err = y.videoDLWorker(destFile, url)
		if err == nil {
			break
		}
	}
	return err
}

func (y *Youtube) StartDownloadWithQuality(destFile string, quality string) error {
	//download highest resolution on [0]
	err := errors.New("Empty stream list")
	for _, v := range y.StreamList {
		if strings.Compare(v["quality"], quality) == 0 {
			url := v["url"]
			y.log(fmt.Sprintln("Download url=", url))

			y.log(fmt.Sprintln("Download to file=", destFile))
			err = y.videoDLWorker(destFile, url)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return y.StartDownload(destFile)
	}
	return err
}

func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

// Values logged with this will only show in debug mode
func (y *Youtube) log(logText string) {
	if y.DebugMode {
		log.Println(logText)
	}
}
