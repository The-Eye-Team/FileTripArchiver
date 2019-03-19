package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func download(worker *sync.WaitGroup, url string) {
	defer worker.Done()

	response, err := doRequest(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		fmt.Println(err)
		return
	}

	downloadPath := ""
	doc.Find(".breadcrumb > div > a > span").Each(func(i int, s *goquery.Selection) {
		downloadPath = downloadPath + "/" + s.Text()
	})

	programName := doc.Find("#content > div.post > h1 > a").First().Text()

	if programName == "" {
		fmt.Println("Cant find name for " + url)
		return
	}

	savePath := arguments.Output + downloadPath + programName

	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := doc.Find(".description").First()
	description := s.Text()

	err = ioutil.WriteFile(savePath+"/description.txt", []byte(description), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	downloadUrl := strings.Replace(url, "download-", "start-", 1)

	response, err = doRequest(downloadUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	var filename string
	var downloadLink string
	scanner := bufio.NewScanner(bytes.NewReader(response))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "software_download_url") {
			r := strings.NewReplacer("software_download_url : '", "", "',", "")
			downloadLink = strings.Trim(r.Replace(line), "	 ")
		}
		if strings.Contains(line, "software_filename") {
			r := strings.NewReplacer("software_filename : '", "", "',", "")
			filename = strings.Trim(r.Replace(line), "	 ")
		}
	}

	if _, err := os.Stat(savePath + "/" + filename); !os.IsNotExist(err) {
		fmt.Println("Skipping File " + filename)
		return
	}

	file, err := os.OpenFile(savePath+"/"+filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}

	response, err = doRequest(downloadLink)
	if err != nil {
		fmt.Println(err)
	}

	_, err = file.Write(response)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done with " + filename)
}

func doRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
