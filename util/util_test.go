package util

import (
	"strings"
	"testing"
)

type ParseTextExpectedIO struct {
	TestCase string
	Result   *ParseTextExpectedResult
}
type ParseTextExpectedResult struct {
	Text    string
	TagList []string
}

func createParseTextExpectedIO(testcase, text string, tagList []string) *ParseTextExpectedIO {
	eResult := &ParseTextExpectedResult{Text: text, TagList: tagList}
	return &ParseTextExpectedIO{TestCase: testcase, Result: eResult}
}

func (self *ParseTextExpectedIO) runTest(t *testing.T) {
	text, tagList := ParseText(self.TestCase)
	if text != self.Result.Text {
		t.Errorf("failed different text(exp: %s <=> act: %s)", text, self.Result.Text)
	}
	tagListJoin := strings.Join(tagList, "")
	selfTagListJoin := strings.Join(self.Result.TagList, "")
	if tagListJoin != selfTagListJoin {
		t.Errorf("failed different tagList(exp: %s <=> act: %s)", tagListJoin, selfTagListJoin)
	}
}

func TestParseText(t *testing.T) {
	input1 := "これはテストです $hoge $fuga"
	expected1 := createParseTextExpectedIO(input1, "これはテストです", []string{"hoge", "fuga"})
	expected1.runTest(t)

	input2 := "これはテストです"
	expected2 := createParseTextExpectedIO(input2, "これはテストです", []string{})
	expected2.runTest(t)

	input3 := "$hoge $fuga これはテストです $goo $boo"
	expected3 := createParseTextExpectedIO(input3, "これはテストです", []string{"hoge", "fuga", "goo", "boo"})
	expected3.runTest(t)

	input4 := "$hoge $fuga これはテストです $hoge $fuga"
	expected4 := createParseTextExpectedIO(input4, "これはテストです", []string{"hoge", "fuga"})
	expected4.runTest(t)
}

var DEFAULT_LOAD_PATH_LIST []string = []string{"../"}

func TestLoadFiles(t *testing.T) {
	result := LoadFiles(DEFAULT_LOAD_PATH_LIST, "config.tml")
	if result != "../config.tml" {
		t.Errorf("failed different filepath(exp: ../config.tml <=> act: %s)", result)
	}
}
