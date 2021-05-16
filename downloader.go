package alphafoxtrot

import (
	"fmt"
	"path"
	"time"

	"github.com/grumpypixel/go-webget"
)

// This isn't the code you're looking for.
// This was created to power the example code.
// And for fun of course.

// Download csv files from OurAirports.com
func DownloadDatabase(targetDir string) {
	files := make([]string, 0)
	for _, filename := range OurAirportsFiles {
		files = append(files, OurAirportsBaseURL+filename)
	}
	for _, url := range files {
		options := webget.Options{
			ProgressHandler: MyProgress{},
			Timeout:         time.Second * 60,
			CreateTargetDir: true,
		}
		webget.DownloadToFile(url, targetDir, "", &options)
	}
}

type MyProgress struct{}

func (p MyProgress) Start(sourceURL string) {
	// fmt.Println()
}

func (p MyProgress) Update(sourceURL string, percentage float64, bytesRead, contentLength int64) {
	fmt.Printf("\rDownloading %s: %v bytes [%.2f%%]", path.Base(sourceURL), bytesRead, percentage)
}

func (p MyProgress) Done(sourceURL string) {
	fmt.Println()
}
