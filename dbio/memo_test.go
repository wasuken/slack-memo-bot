package dbio

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wasuken/slack-memo-bot/util"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const TEST_DB_PATH = "./test.sqlite"

func setup() *sql.DB {
	_, err := os.Stat(TEST_DB_PATH)
	if err == nil {
		os.Remove(TEST_DB_PATH)
	}
	os.Create(TEST_DB_PATH)
	f, err := os.Open("../create.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", TEST_DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(string(b))
	return db
}

func TestInsertTags(t *testing.T) {
	db := setup()
	defer db.Close()
	tagList := []string{"hoge", "fuga"}
	insertTags(db, tagList)

	actual := len(selectInsertedTagList(db))
	expected := len(tagList)
	if actual != expected {
		t.Errorf("different rec count %d <=> %d.", actual, expected)
	}
	os.Remove(TEST_DB_PATH)
}

func TestSaveMemo(t *testing.T) {
	db := setup()
	defer db.Close()

	text, tagList := util.ParseText("これはテスト $hoge $fuga")
	SaveMemo(db, text, tagList)

	// tags check
	expectedTagList := []string{"hoge", "fuga"}
	tl := tagsListFromDB(db)
	for _, recAry := range tl {
		if !util.Contains(expectedTagList, recAry[1]) {
			t.Errorf("not contains tags(%s in %s).", recAry[1], expectedTagList)
		}
	}
	// memos check
	expectedMemoList := []string{"これはテスト"}
	ml := memosListFromDB(db)
	for _, rec := range ml {
		if !util.Contains(expectedMemoList, rec[1]) {
			t.Errorf("failed insert memo contents(%s).", rec[1])
		}
	}

	// memo_tags check
	expectedMemoTagsList := [][]string{{"1", "1"}, {"1", "2"}}
	mtl := memoTagsListFromDB(db)
	for _, rec := range mtl {
		if !AryContains(expectedMemoTagsList, rec) {
			t.Errorf("failed inserted memo_tags(memo_id:%s, tag_id:%s)", rec[0], rec[1])
		}
	}

	os.Remove(TEST_DB_PATH)
}

// genericsが使えるバージョンなら...
func AryContains(haystack [][]string, needle []string) bool {
	needleConcat := strings.Join(needle, "")
	for _, ary := range haystack {
		if strings.Join(ary, "") == needleConcat {
			return true
		}
	}
	return false
}

// 型変換だるいのでレコードは強制的にstringに変換される
func tagsListFromDB(db *sql.DB) (result [][]string) {
	result = [][]string{}
	rows := queryRows(db, "select * from tags")
	defer rows.Close()

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		result = append(result, []string{id, name})
	}
	return result
}

// 型変換だるいのでレコードは強制的にstringに変換される
func memosListFromDB(db *sql.DB) (result [][]string) {
	result = [][]string{}
	rows := queryRows(db, "select id, contents from memos")
	defer rows.Close()

	for rows.Next() {
		var id, contents string
		if err := rows.Scan(&id, &contents); err != nil {
			log.Fatal(err)
		}
		result = append(result, []string{id, contents})
	}
	return result
}

// 型変換だるいのでレコードは強制的にstringに変換される
func memoTagsListFromDB(db *sql.DB) (result [][]string) {
	result = [][]string{}
	rows := queryRows(db, "select memo_id, tag_id from memo_tags")
	defer rows.Close()

	for rows.Next() {
		var memo_id, tag_id string
		if err := rows.Scan(&memo_id, &tag_id); err != nil {
			log.Fatal(err)
		}
		result = append(result, []string{memo_id, tag_id})
	}
	return result
}

func queryRows(db *sql.DB, query string) *sql.Rows {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
