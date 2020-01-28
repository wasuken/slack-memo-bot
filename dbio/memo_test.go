package dbio

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

// func TestParseTextSuccess(t *testing.T) {
// 	text, tagList := parseText("$hoge $fuga うんちっち")
// 	if text != "うんちっち" {
// 		t.Errorf("got: %v\nwant: %v", text, "うんちっち")
// 	}
// 	actual := []string{"hoge", "fuga"}
// 	if len(tagList) != len(actual) {
// 		t.Errorf("got: %v\nwant: %v", tagList, []string{"hoge", "fuga"})
// 	}
// 	for i := 0; i < len(actual); i++ {
// 		if actual[i] != tagList[i] {
// 			t.Errorf("got: %v\nwant: %v", tagList, []string{"hoge", "fuga"})
// 		}
// 	}
// }
