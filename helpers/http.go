package helpers

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetFile(url string) *os.File {
	tmpfile, err := ioutil.TempFile("", "template")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(tmpfile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n)
	return tmpfile
}
