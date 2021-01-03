package sqlite

import (
	"database/sql"
	"crypto/rand"
	"github.com/Fliegermarzipan/gallipot/data"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	saltFetchQuery = "SELECT salt FROM user WHERE username = ?"
	sessionCreateQuery = "INSERT INTO session (userid) "+
		"SELECT userid FROM user WHERE username=? AND salt=? AND hash=?"
	sessionFromIdQuery = "SELECT username, displayname, email, token, last_access, expires "+
		"FROM user JOIN session WHERE session.rowid=?"
	userRegisterQuery = "INSERT INTO user(username, email, salt, hash) VALUES(?,?,?,?)"
	userFromIdQuery = "SELECT username, displayname, email FROM user WHERE rowid=?"
)

type SQLitePersist struct {
	db *sql.DB
	saltFetchStmt *sql.Stmt
	sessionCreateStmt *sql.Stmt
	sessionFromIdStmt *sql.Stmt
	userRegisterStmt *sql.Stmt
	userFromIdStmt *sql.Stmt
}

func NewSQLitePersist(filename string) (persist *SQLitePersist) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	persist = new(SQLitePersist)

	persist.db = db
	return persist
}

func (p *SQLitePersist) getSalt(tx *sql.Tx, ld data.LoginData) (salt []byte) {
	if p.saltFetchStmt == nil {
		var err error
		p.saltFetchStmt, err = p.db.Prepare(saltFetchQuery)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Stmt(p.saltFetchStmt).QueryRow(ld.Username).Scan(&salt); err != nil {
		log.Fatal(err)
	}
	return
}

func (p *SQLitePersist) Authenticate(ld data.LoginData) *data.Session {
	tx,err := p.db.Begin()
	if err != nil{
		log.Fatal(err)
	}
	salt := p.getSalt(tx, ld)
	

	if p.sessionCreateStmt == nil{
		p.sessionCreateStmt, err = p.db.Prepare(sessionCreateQuery)
		if err != nil {
			log.Fatal(err)
		}
	}

	rs, err := tx.Stmt(p.sessionCreateStmt).Exec(ld.Username, salt, ld.Hash(salt))
	if err != nil{
		log.Fatal(err)
	}

	var session data.Session
	session.User = new(data.User)

	if p.sessionFromIdStmt == nil{
		p.sessionFromIdStmt, err = p.db.Prepare(sessionFromIdQuery)
		if err != nil{
			log.Fatal(err)
		}
	}

	rowid, err := rs.LastInsertId()
	var lastActive, expiry int64
	err = tx.Stmt(p.sessionFromIdStmt).QueryRow(rowid).Scan(&(session.User.Username),
		&(session.User.Displayname), &(session.User.Email), &(session.Token),
		&lastActive, &expiry)
	session.LastActive = time.Unix(lastActive,0)
	session.Expiry = time.Unix(lastActive,0)
	if err != nil{
		log.Fatal(err)
	}
	tx.Commit()
	return &session
}



func (p *SQLitePersist) Register(ld data.LoginData, email string)(u *data.User){
	u = new(data.User)
	salt := make([]byte, 36)
	_, err := rand.Read(salt)
	if err!=nil{
		log.Fatal(err)
	}
	hash := ld.Hash(salt)
	tx, err := p.db.Begin()
	if err!=nil{
		log.Fatal(err)
	}
	if p.userRegisterStmt == nil{
		p.userRegisterStmt,err = p.db.Prepare(userRegisterQuery)
		if err!=nil{
			log.Fatal(err)
		}
	}
	rs, err := tx.Stmt(p.userRegisterStmt).Exec(ld.Username, email, salt, hash)
	if err!=nil{
		log.Fatal(err)
	}
	rowid, err := rs.LastInsertId()
	if err!=nil{
		log.Fatal(err)
	}

	if p.userFromIdStmt == nil{
		p.userFromIdStmt,err = p.db.Prepare(userFromIdQuery)
		if err!=nil{
			log.Fatal(err)
		}
	}
	err = tx.Stmt(p.userFromIdStmt).QueryRow(rowid).Scan(&(u.Username), &(u.Displayname), &(u.Email))
	if err!=nil{
		log.Fatal(err)
	}
	tx.Commit()
	return
}
