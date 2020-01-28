create table memos(
	   id integer primary key,
	   contents text,
	   created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	   updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);
create table tags(
	   id integer primary key,
	   name text
);
create table memo_tags(
	   memo_id integer,
	   tag_id integer,
	   primary key(memo_id, tag_id)
);
