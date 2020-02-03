package util

import (
	"os"
	"os/user"
	"strings"
)

func ParseText(text string) (string, []string) {
	nodes := strings.Split(text, " ")
	textList := []string{}
	tagList := []string{}
	for _, node := range nodes {
		if node[0:1] == "$" {
			rep := strings.Replace(node, "$", "", -1)
			if !Contains(tagList, rep) {
				tagList = append(tagList, rep)
			}
		} else {
			textList = append(textList, node)
		}
	}
	return strings.Join(textList, ""), tagList
}

func Contains(lst []string, s string) bool {
	for _, v := range lst {
		if v == s {
			return true
		}
	}
	return false
}

func LoadFiles(filepaths []string, filename string) string {
	usr, _ := user.Current()
	f := strings.Replace(filepaths[0], "~", usr.HomeDir, 1) + filename
	_, err := os.Stat(f)
	if err != nil {
		return LoadFiles(filepaths[1:], filename)
	} else {
		return f
	}
}
