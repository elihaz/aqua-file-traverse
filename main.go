package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const numWorkers = 4

var root string
var address string

func init() {

	flag.StringVar(&root, "path", "", "traverse root path")
	flag.StringVar(&address, "address", "/", "HTTP service address")
	flag.Parse()

}

func main() {

	e := validate(address, root)
	if e != nil {
		log.Println(e)
		return
	}

	done := make(chan bool)
	infoChan := make(chan *FileInformation)

	log.Print(fmt.Sprintf("starting traverse: root path: %s, will send post request to: %s",
		root, address))

	go traverse(root, infoChan)

	client := &http.Client{}

	for i := 0; i < numWorkers; i++ {

		go func(fileInfo chan *FileInformation) {
			for {
				j, more := <-fileInfo
				if more {
					marshal, err := json.Marshal(j)
					if err != nil {
						log.Printf("json marshal failed: %v", err)
						continue
					}
					log.Print(j)
					e := sendInfoData(address, marshal, client)
					if e != nil {
						log.Println("error sending data: " + e.Error())
					}

				} else {
					done <- true
				}
			}
		}(infoChan)
	}
	<-done
	log.Print("<------ finished traversing succesfully ---------->")
}

func validate(address string, root string) error {

	_, parseUrlErr := url.ParseRequestURI(address)
	if parseUrlErr != nil {
		return parseUrlErr
	}

	_, pathErr := exists(root)
	if pathErr != nil {
		return pathErr
	}

	return nil
}

func sendInfoData(url string, jsonStr []byte, client *http.Client) error {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response, e := client.Do(req)

	if e != nil {
		return fmt.Errorf("got error posting to server: %s ", e)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("post didn't succeed. response: %s ", response)
	}

	return nil

}

func traverse(root string, infoChan chan *FileInformation) error {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("os.Stat failed: %s ", err)
		}

		information := &FileInformation{
			Name: info.Name(),
			Size: info.Size(),
			Ext:  filepath.Ext(info.Name()),
		}

		mode := fi.Mode()
		if mode.IsRegular() {
			infoChan <- information
		}

		return nil
	})

	close(infoChan)

	return err
}

type FileInformation struct {
	Name string
	Size int64
	Ext  string
}

func exists(path string) (bool, error) {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
