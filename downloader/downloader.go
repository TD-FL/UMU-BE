package downloader

import (
	"github.com/yeka/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var zipPassword = "zaaeF223hfkdfsDF243DFGSfsdfSDFAsdfadfafSDFASDFASFppOIjSA343423Lrfus34324325gysdx3v"

func Download(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Extract(zipPath string) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(zipPassword)
		}
		r, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		if f.FileInfo().IsDir() {
			os.MkdirAll("./urnik", os.ModePerm)
		} else {
			ioutil.WriteFile("./timetable.csv", buf, os.ModePerm)
		}
	}
}
