package util

import (
	"strings"
)

func ParseText(text string) (string, []string) {
	nodes := strings.Split(text, " ")
	textList := []string{}
	tagList := []string{}
	for _, node := range nodes {
		if node[0:1] == "$" {
			tagList = append(tagList, strings.Replace(node, "$", "", -1))
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
