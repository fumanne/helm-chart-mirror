package main

import (
	"gopkg.in/yaml.v2"
	"helm-chart-mirror/chart"
	"helm-chart-mirror/fetch"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	OFFICAL_URL = "https://kubernetes-charts.storage.googleapis.com/index.yaml"
	MAXFILE     = 2048
)

func main() {
	var wg sync.WaitGroup
	maxChan := make(chan bool, MAXFILE)
	I := chart.Index{}
	if err := yaml.Unmarshal(fetch.FetchIndexYaml(OFFICAL_URL), &I); err != nil {
		log.Fatalf("Unmarshal struct %s Failed, Error is %s\n", I, err)
	}
	for _, charts := range I.Entries {
		for _, chart := range charts {
			wg.Add(1)
			maxChan <- true
			go chart.Download(&wg, maxChan)
		}
	}
	wg.Wait()
	I.SetGenerated()
	d, err := yaml.Marshal(&I)
	if err != nil {
		log.Fatalf("Marshal struct %s Failed, Error is %s\n", I, err)
	}
	if err := ioutil.WriteFile("./docs/index.yaml", d, os.ModePerm); err != nil {
		log.Fatalf("Write Index.yaml Failed, Error is %s\n", err)
	}

}
