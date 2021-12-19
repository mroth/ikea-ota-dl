package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	feedStable = "http://fw.ota.homesmart.ikea.net/feed/version_info.json"
	feedTest   = "http://fw.test.ota.homesmart.ikea.net/feed/version_info.json"
)

var (
	verbose    = flag.Bool("v", true, "verbose")
	feedURL    = flag.String("feed", feedStable, "firmware update feed URL")
	numWorkers = flag.Int("workers", 4, "max concurrent downloads")
)

type FeedEntry struct {
	BinaryURL        string `json:"fw_binary_url"`
	FileVersionLSB   int    `json:"fw_file_version_LSB,omitempty"`
	FileVersionMSB   int    `json:"fw_file_version_MSB,omitempty"`
	Filesize         int    `json:"fw_filesize"`
	ImageType        int    `json:"fw_image_type,omitempty"`
	ManufacturerID   int    `json:"fw_manufacturer_id,omitempty"`
	Type             int    `json:"fw_type"`
	HotfixVersion    int    `json:"fw_hotfix_version,omitempty"`
	MajorVersion     int    `json:"fw_major_version,omitempty"`
	MinorVersion     int    `json:"fw_minor_version,omitempty"`
	ReqHotfixVersion int    `json:"fw_req_hotfix_version,omitempty"`
	ReqMajorVersion  int    `json:"fw_req_major_version,omitempty"`
	ReqMinorVersion  int    `json:"fw_req_minor_version,omitempty"`
	UpdatePrio       int    `json:"fw_update_prio,omitempty"`
	WeblinkRelnote   string `json:"fw_weblink_relnote,omitempty"`
	BuildVersion     int    `json:"fw_build_version,omitempty"`
}

// Filename returns the filename portion of a FeedEntry.BinaryURL.
func (fe FeedEntry) Filename() string {
	return filepath.Base(fe.BinaryURL)
}

// ParseFeed retrieves and parses an IKEA firmware feed for given url.
func ParseFeed(url string) ([]FeedEntry, error) {
	var results []FeedEntry
	r, err := http.Get(url)
	if err != nil {
		return results, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&results)
	return results, err
}

// DownloadFile will download url to a local file.
func DownloadFile(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666) // err if exists
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <destination>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	destination := args[0]

	// parse feed
	entries, err := ParseFeed(*feedURL)
	if err != nil {
		log.Fatal(err)
	}

	// queue up feed entries for retrieval
	ch := make(chan FeedEntry)
	go func() {
		for _, e := range entries {
			ch <- e
		}
		close(ch)
	}()

	// spawn retrieval downloaders
	var wg sync.WaitGroup
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go func() {
			for e := range ch {
				dstpath := filepath.Join(destination, e.Filename())
				err := DownloadFile(dstpath, e.BinaryURL)
				if errors.Is(err, os.ErrExist) && *verbose {
					log.Println(dstpath, "already exists")
				} else if err == nil && *verbose {
					log.Println(dstpath, "downloaded")
				} else {
					log.Println(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
