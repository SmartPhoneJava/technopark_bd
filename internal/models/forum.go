package models

type Forum struct {
	Posts   int    `json:"-" db:"posts"`
	Threads int    `json:"-" db:"threads"`
	Slug    string `json:"slug" db:"slug"`
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user_nickname"`
}

/*
 CREATE Table Forum (
        posts int default 0,
        slug SERIAL PRIMARY KEY,
        threads int,
        title varchar(60) not null,
        user_nickname varchar(80) not null
    );
*/
