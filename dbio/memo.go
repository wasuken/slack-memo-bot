package dbio

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wasuken/slack-memo-bot/util"
	"log"
)

func DeleteMemoTags(db *sql.DB, tagList []string) {
	for _, tag := range tagList {
		deleteMemoTag(db, tag)
	}
}

func deleteMemoTag(db *sql.DB, tag string) {
	rows, err := db.Query("SELECT memos.id as mid, name FROM memos join memo_tags as mt on mt.memo_id = memos.id join tags on mt.tag_id = tags.id where name = ?", tag)
	if err != nil {
		log.Fatal(err)
	}
	midList := []int{}
	nameList := []string{}
	for rows.Next() {
		var mid int
		var name string
		if err := rows.Scan(&mid, &name); err != nil {
			log.Fatal(err)
		}
		midList = append(midList, mid)
		nameList = append(nameList, name)
	}
	rows.Close()
	for i := 0; i < len(midList); i++ {
		_, err := db.Exec("delete from memo_tags where memo_id = ?", midList[i])
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("delete from memos where id = ?", midList[i])
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("delete from tags where name = ?", nameList[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// メモを全てまとめ、結合した文字列を返す。
func OutputMemo(db *sql.DB, outputType string, tagList []string, user string) (message string) {
	// tagごとに並び替えを行い、その中で生成時間順で並び替える。
	rows, err := db.Query("SELECT name, contents, updated_at FROM memos join memo_tags as mt on mt.memo_id = memos.id join tags on mt.tag_id = tags.id where user = ? order by name desc, date(updated_at, 'localtime')", user)
	tagMemoMap := map[string]string{}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var contents string
		var name string
		var updated_at string
		if err := rows.Scan(&name, &contents, &updated_at); err != nil {
			log.Fatal(err)
		}
		tagMemoMap[name] += "\n" + contents + "\n"
	}
	switch outputType {
	case "markdown":
		for k, v := range tagMemoMap {
			message += "# " + k + "\n\n"
			message += v + "\n\n"
		}
	}
	return message
}

func SaveMemo(db *sql.DB, text string, tagList []string, user string) string {
	// insert tags table, and get tags_id.
	tagrecs := insertTags(db, tagList)
	// insert to memo table, and get memo_id.
	memoId := insertMemo(db, text, user)
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

func insertMemo(db *sql.DB, text, user string) int64 {
	result, err := db.Exec("insert into memos(contents, user) values(?, ?)", text, user)
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
		if !util.Contains(inserted_tag_list, tag) {
			_, err := db.Exec("insert into tags(name) values(?)", tag)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	inserted_tagrec_list = selectInsertedTagList(db)
	inserting_tagrec_list := []*TagRecord{}
	for _, tagrec := range inserted_tagrec_list {
		if util.Contains(tagList, tagrec.name) {
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
