package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const mainPage = "https://www.istockphoto.com/"
const toadQuery = "ru/search/2/film?phrase=common+toad"

func main() {
	toadLinks, err := parseMainPage()
	if err != nil {
		log.Println("[Parser] Error while parsing main page")
		return
	}

	log.Println(len(toadLinks))

	for id, link := range toadLinks {
		parseToadPage(link, id)
	}

}

func parseMainPage() (result []string, err error) {
	responce, err := http.Get(mainPage + toadQuery)
	if err != nil {
		log.Println("[Parser] Error while parsing main page")
		return result, err
	}

	body, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		log.Println("[Parser] Error while reading responce.Body")
		return result, err
	}

	regular := regexp.MustCompile("\"\\/%.{100,600}\\d{9}-\\d{9}")
	finded := regular.FindAll(body, -1)

	for _, link := range finded {
		result = append(result, mainPage+string(link)[2:])
	}
	return result, err
}

func parseToadPage(toadLink string, toadNum int) error {
	log.Println("next")
	responce, err := http.Get(toadLink)
	if err != nil {
		log.Println("[Parser] Error while parsing main page")
		return err
	}

	body, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		log.Println("[Parser] Error while reading responce.Body")
		return err
	}

	regular := regexp.MustCompile("https://media.istockphoto.com.[^j\\\\]*?=\"")
	finded := regular.FindAll(body, 1)
	var videoLink string
	if len(finded) > 0 {
		videoLink = string(finded[0])
		videoLink = strings.ReplaceAll(videoLink, "amp;", "")
		videoLink = videoLink[:len(videoLink)-1]
		err := downloadVideo(videoLink, toadNum)
		if err != nil {
			log.Println("[Parser] error to download file: ", err)
			return err
		}
	}
	return nil
}

func downloadVideo(link string, toadId int) (err error) {
	file, err := os.Create("toads_and_frogs/" + fmt.Sprint(toadId) + ".mp4")
	if err != nil {
		log.Println("[Downloader] Failed to create file: ", err)
		return err
	}

	responce, err := http.Get(link)
	if err != nil {
		log.Println("[Downloader] Failed to get http responce: ", err)
		return err
	}

	_, err = io.Copy(file, responce.Body)
	if err != nil {
		log.Println("[Downloader] Failed to copy to file: ", err)
		return err
	}
	return nil
}
