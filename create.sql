create table memos(
	   id integer primary key,
	   contents text
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
