package chart

import (
	"os"
	"testing"
)

func TestFetchGitEnv(t *testing.T) {
	os.Setenv("GIT_USER", "FU")
	s := FetchGitEnv("GIT_USER", "HAHA")
	if s == "HAHA" {
		t.Fatalf("Get Env key GIT_USER is not expected %s", s)
	}
}
