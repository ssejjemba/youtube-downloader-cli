package provider

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func videoDLWorker(destFile string, target string) error {
	resp, err := http.Get(target)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("reading answer: non 200 status code received: '%s'", err)
		return errors.New("non 200 status code received")
	}

	out, err := os.Create(destFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("download video err=", err)
		return err
	}
	return nil
}