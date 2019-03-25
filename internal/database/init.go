package database

import (
	"database/sql"
	"escapade/internal/config"
	"fmt"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, "user=rolepade password=escapade dbname=escabase sslmode=disable")
	}

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		fmt.Println("database/Init cant open:" + err.Error())
		return
	}

	db = &DataBase{
		Db: database,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		fmt.Println("database/Init cant access:" + err.Error())
		return
	}
	fmt.Println("database/Init open")
	//if !db.areTablesCreated(CDB.Tables) {
	if err = db.CreateTables(); err != nil {
		return
	}
	//}

	return
}

func (db *DataBase) checkTable(tableName string) (err error) {
	sqlStatement := `
    SELECT count(1)
  FROM information_schema.tables tbl 
  where tbl.table_name like $1;`
	row := db.Db.QueryRow(sqlStatement, tableName)

	var result int
	if err = row.Scan(&result); err != nil {
		fmt.Println(tableName + " doesnt exists. Create it!" + err.Error())

		return
	}
	return
}

func (db *DataBase) areTablesCreated(tables []string) (created bool) {
	created = true
	for _, table := range tables {
		if err := db.checkTable(table); err != nil {
			created = false
			break
		}
	}
	return
}

func (db *DataBase) CreateTables() error {
	sqlStatement := `

    DROP TABLE IF EXISTS Vote;
    DROP TABLE IF EXISTS Post;
    DROP TABLE IF EXISTS Thread;
    DROP TABLE IF EXISTS Forum;
    DROP TABLE IF EXISTS UserForum;

    CREATE Table UserForum (
        id SERIAL PRIMARY KEY,
        nickname varchar(80) UNIQUE NOT NULL,
        fullname varchar(30) NOT NULL,
        email varchar(50) UNIQUE NOT NULL,
        about varchar(1000) 
    );

    CREATE Table Forum (
        id SERIAL PRIMARY KEY,
        posts int default 0,
        slug varchar(80) not null UNIQUE,
        threads int default 0,
        title varchar(120) not null,
        user_nickname varchar(80) not null
    );

    ALTER TABLE Forum
        ADD CONSTRAINT forum_user
        FOREIGN KEY (user_nickname)
        REFERENCES UserForum(nickname)
            ON DELETE CASCADE;

    CREATE Table Thread (
        id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        forum varchar(120) not null,
        message varchar(2100) not null,
        created    TIMESTAMPTZ,
        title varchar(120) not null,
        votes int default 0,
        slug varchar(120) default null
    );

    ALTER TABLE Thread
    ADD CONSTRAINT thread_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Thread
    ADD CONSTRAINT thread_forum
    FOREIGN KEY (forum)
    REFERENCES Forum(slug)
        ON DELETE CASCADE;

    CREATE Table Post (
        id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        forum varchar(120),
        message varchar(2400) not null,
        created    TIMESTAMPTZ,
        isEdited boolean default false,
        thread int,
        parent int,
        path varchar(1000),
        level int
    );

    ALTER TABLE Post
    ADD CONSTRAINT post_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Post
    ADD CONSTRAINT post_forum
    FOREIGN KEY (forum)
    REFERENCES Forum(slug)
        ON DELETE CASCADE;
    
    ALTER TABLE Post
    ADD CONSTRAINT post_thread
    FOREIGN KEY (thread)
    REFERENCES Thread(id)
        ON DELETE CASCADE;

    CREATE Table Vote (
        id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        thread int not null,
        isEdited bool default false,
        voice int default 0
    );

    ALTER TABLE Vote
    ADD CONSTRAINT vote_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Vote
    ADD CONSTRAINT vote_thread
    FOREIGN KEY (thread)
    REFERENCES Thread(id)
        ON DELETE CASCADE;

	`
	_, err := db.Db.Exec(sqlStatement)

	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}
	return err
}
