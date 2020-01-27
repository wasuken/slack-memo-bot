package dbio

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func SaveMemo(slackText string) string {
	text, tagList := parseText(slackText)

	if len(tagList) <= 0 {
		return "no tags"
	}

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// insert tags table, and get tags_id.
	tagrecs := insertTags(db, tagList)
	// insert to memo table, and get memo_id.
	memoId := insertMemo(db, text)
	// insert to memo_tags table from inserted tags and memo
	insertTagMemos(db, tagrecs, memoId)
	return "save!"
}

func insertTagMemos(db *sql.DB, tagrecs []*TagRecord, memoId int64) {
	for _, tagrec := range tagrecs {
		_, err := db.Exec("insert into memo_tags(memo_id, tag_id) values(?, ?)", memoId, tagrec.id)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertMemo(db *sql.DB, text string) int64 {
	result, err := db.Exec("insert into memos(contents) values(?)", text)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}
func insertTags(db *sql.DB, tagList []string) []*TagRecord {
	inserted_tagrec_list := selectInsertedTagList(db)
	inserted_tag_list := []string{}
	for _, tagrec := range inserted_tagrec_list {
		inserted_tag_list = append(inserted_tag_list, tagrec.name)
	}
	for _, tag := range tagList {
		if !contains(inserted_tag_list, tag) {
			_, err := db.Exec("insert into tags(name) values(?)", tag)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	inserted_tagrec_list = selectInsertedTagList(db)
	inserting_tagrec_list := []*TagRecord{}
	for _, tagrec := range inserted_tagrec_list {
		if contains(tagList, tagrec.name) {
			inserting_tagrec_list = append(inserting_tagrec_list, tagrec)
		}
	}
	return inserting_tagrec_list
}

type TagRecord struct {
	id   int
	name string
}

func createTagRecord(id int, name string) *TagRecord {
	return &TagRecord{name: name, id: id}
}

func selectInsertedTagList(db *sql.DB) []*TagRecord {
	inserted_tagrec_list := []*TagRecord{}
	rows, err := db.Query("SELECT id, name FROM tags")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		inserted_tagrec_list = append(inserted_tagrec_list, createTagRecord(id, name))
	}
	return inserted_tagrec_list
}

func contains(lst []string, s string) bool {
	for _, v := range lst {
		if v == s {
			return true
		}
	}
	return false
}

func parseText(text string) (string, []string) {
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
