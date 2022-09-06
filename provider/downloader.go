package provider

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	if(y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
		y.DownloadPercentage <- int64(y.downloadLevel)
	}
	return
}

func (y *Youtube) videoDLWorker(destFile string, target string) error {
	resp, err := http.Get(target)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
		return err
	}

	defer resp.Body.Close()
	y.contentLength = float64(resp.ContentLength)

	if resp.StatusCode != 200 {
		log.Printf("reading answer: non 200[code=%v] status code received: '%v'", resp.StatusCode, err)
		return errors.New("non 200 status code received")
	}

	// create dir structure if none exists
	err = os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(destFile)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(out, y)
	_, err = io.Copy(mw, resp.Body)
	if err != nil {
		log.Println("download video err=", err)
		return err
	}
	return nil
}