package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Title   string   `yaml:"title"`
	PubDate string   `yaml:"PubDate"`
	Tags    []string `yaml:"tags"`
	Author  string   `yaml:"author"`
}

func getParts(filename string) (parts [][]byte, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parts = bytes.SplitN(content, []byte("---"), 3)

	return parts, nil
}
func readMetadata(meta []byte) (metadata Metadata, err error) {
	var metaData Metadata

	if err := yaml.Unmarshal(meta, &metaData); err != nil {
		return Metadata{}, err
	}

	return metaData, nil

}

func main() {
	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		rss := generateFeed()
		w.Write([]byte(rss))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("error starting server", err)
		os.Exit(1)
	}
}
func generateFeed() (rss string) {
	const dir = "./content"

	feed := &feeds.Feed{
		Title:       "DevOps Weekly Updates",
		Link:        &feeds.Link{Href: "http://devops.elliottmgmt.com"},
		Description: "An RSS Feed of Elliott's DevOps team updates",
		Created:     time.Now(),
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading content directory")
		os.Exit(1)
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) == ".md" {
			parts, err := getParts(filepath.Join(dir, file.Name()))

			if err != nil {
				fmt.Println("Error processing file: ", err)
				continue
			}
			metadata, err := readMetadata(parts[1])
			html := blackfriday.Run(parts[2])
			fmt.Println("html: ", string(html), len(parts))

			pubDate, err := time.Parse("2006-01-02", metadata.PubDate)
			item := &feeds.Item{
				Title:       metadata.Title,
				Link:        &feeds.Link{Href: "http://devops.elliottmgmt.com/" + file.Name()},
				Description: string(html),
				Created:     pubDate,
			}

			feed.Add(item)
		}

	}
	rss, err = feed.ToRss()
	if err != nil {
		fmt.Println("Error generating feed:", err)
		os.Exit(1)
	}

	fmt.Println(rss)
	return (rss)
}
