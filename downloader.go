package alphafoxtrot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// This isn't the code you're looking for.
// This was created to power the example code.
// And for fun of course.

const tempExtension = ".download"

// Download csv files from OurAirports.com
func DownloadDatabase(targetDir string) {
	files := make([]string, 0)
	for _, filename := range OurAirportsFiles {
		files = append(files, OurAirportsBaseURL+filename)
	}
	for _, url := range files {
		filename := path.Base(url)
		target := path.Join(targetDir, filename)
		temp := target + tempExtension
		label := filename
		if err := downloadFile(temp, url, label); err != nil {
			log.Println(err)
			continue
		}
		if err := os.Rename(temp, target); err != nil {
			log.Println(err)
		}
		fmt.Println("")
	}
}

func downloadFile(filepath, url, progressLabel string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	progress := &FileProgress{Name: progressLabel, Size: uint64(resp.ContentLength)}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, progress)); err != nil {
		return err
	}
	return nil
}

type FileProgress struct {
	Name  string
	Bytes uint64
	Size  uint64
}

func (fp *FileProgress) Write(p []byte) (int, error) {
	n := len(p)
	fp.Bytes += uint64(n)
	fp.UpdateProgress()
	return n, nil
}

func (fp *FileProgress) UpdateProgress() {
	percentage := ""
	if fp.Size > 0 {
		value := float64(fp.Bytes) / float64(fp.Size) * 100.0
		percentage = fmt.Sprintf("[%.0f%%]", value)
	}
	fmt.Printf("\rDownloading %s: %v bytes %s", fp.Name, fp.Bytes, percentage)
}
