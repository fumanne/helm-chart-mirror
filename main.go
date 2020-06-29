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

const OFFICAL_URL = "https://kubernetes-charts.storage.googleapis.com/index.yaml"

func main() {
	var wg sync.WaitGroup

	I := chart.Index{}
	if err := yaml.Unmarshal(fetch.FetchIndexYaml(OFFICAL_URL), &I); err != nil {
		log.Fatalf("Unmarshal struct %s Failed, Error is %s\n", I, err)
	}
	for _, charts := range I.Entries {
		for _, chart := range charts {
			wg.Add(1)
			go chart.Download(&wg)
		}
	}
	wg.Wait()
	d, err := yaml.Marshal(&I)
	if err != nil {
		log.Fatalf("Marshal struct %s Failed, Error is %s\n", I, err)
	}
	if err := ioutil.WriteFile("./index.yaml", d, os.ModePerm); err != nil {
		log.Fatalf("Write Index.yaml Failed, Error is %s\n", err)
	}

}
