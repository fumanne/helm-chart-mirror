package fetch

import (
	"github.com/levigross/grequests"
	"log"
)

func FetchIndexYaml(url string) []byte {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		log.Fatalf("Download index.yaml Failed, Error is %s\n", err)
	}
	return resp.Bytes()
}
