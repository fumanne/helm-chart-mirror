package chart

import "os"

func FetchGitEnv(name, defaultvalue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = defaultvalue
	}
	return value
}
