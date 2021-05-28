package sqlite

import (
	"crypto/rand"
	"database/sql"
	_ "embed"
	"github.com/TMDNCM/ictm/data"
	"github.com/TMDNCM/ictm/persistence"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	saltFetchQuery     = "SELECT salt FROM user WHERE username = ?"
	sessionCreateQuery = "INSERT INTO session (userid) " +
		"SELECT userid FROM user WHERE username=? AND salt=? AND hash=?"
	sessionFromIdQuery = "SELECT username, displayname, email, token, last_access, expires " +
		"FROM user NATURAL JOIN session WHERE session.sessionid=?"
	sessionIdFromTokenQuery = "SELECT sessionid FROM session WHERE token = ?"
	sessionExpiryQuery      = "SELECT expires FROM session WHERE sessionid = ?"
	sessionInvalidateQuery  = "DELETE FROM session WHERE sessionid=?"
	userRegisterQuery       = "INSERT INTO user(username, email, salt, hash) VALUES(?,?,?,?)"
	userFromIdQuery         = "SELECT username, displayname, email FROM user WHERE userid=?"
	userIdFromNameQuery     = "SELECT userid FROM user WHERE username = ?"
	setUsernameQuery        = "UPDATE user SET username = ? WHERE userid=?"
	setEmailQuery           = "UPDATE user SET email = ? WHERE userid=?"
	setDisplaynameQuery     = "UPDATE user SET displayname = ? WHERE userid=?"
	setHashQuery            = "UPDATE user SET hash =? WHERE userid=? "
	getFriendsQuery         = "SELECT friend FROM friends WHERE user = ? AND confirmed = 1"
	getHistoryQuery         = "SELECT entryid FROM logentry WHERE user = ?"
	logDoseQuery            = "INSERT INTO logentry(user, substance, route, amount, unit, taken) VALUES(?,?,?,?,?,?)"
	addFriendQuery          = "INSERT INTO friends_t (user1, user2) VALUES(?,?)"
	confirmFriendQuery      = "UPDATE friends_t SET comfirmed=1 WHERE user2=? AND user1=?"
	doseSetWhenQuery        = "UPDATE logentry SET taken = ? WHERE entryid=?"
	doseSetAmountQuery      = "UPDATE logentry set amount = ? WHERE entryid=?"
	doseSetSubstanceQuery   = "UPDATE logentry SET substance = ? WHERE entryid = ?"
	doseSetRouteQuery       = "UPDATE logentry SET route = ? WHERE entryid = ?"
	getDoseQuery            = "SELECT user, substance, route, amount, unit, taken FROM logentry WHERE entryid=?"
)

//go:embed proto.sql
var initQuery string

type SQLitePersist struct {
	db                     *sql.DB
	saltFetchStmt          *sql.Stmt
	sessionCreateStmt      *sql.Stmt
	sessionFromIdStmt      *sql.Stmt
	sessionIdFromTokenStmt *sql.Stmt
	sessionExpiryStmt      *sql.Stmt
	sessionInvalidateStmt  *sql.Stmt
	userRegisterStmt       *sql.Stmt
	userFromIdStmt         *sql.Stmt
	userIdFromNameStmt     *sql.Stmt
	setUsernameStmt        *sql.Stmt
	setEmailStmt           *sql.Stmt
	setDisplaynameStmt     *sql.Stmt
	setHashStmt            *sql.Stmt
	getFriendsStmt         *sql.Stmt
	getHistoryStmt         *sql.Stmt
	logDoseStmt            *sql.Stmt
	addFriendStmt          *sql.Stmt
	doseSetWhenStmt        *sql.Stmt
	doseSetAmountStmt      *sql.Stmt
	doseSetSubstanceStmt   *sql.Stmt
	doseSetRouteStmt       *sql.Stmt
	getDoseStmt            *sql.Stmt
	confirmFriendStmt      *sql.Stmt
}

type Session struct {
	*SQLitePersist
	sessionid uint64
}

type User struct {
	*SQLitePersist
	userid uint64
}

type Doses struct {
	*SQLitePersist
	query string
	args  []interface{}
}

type Dose struct {
	*SQLitePersist
	doseid uint64
}

func NewPersistor(filename string) (persist *SQLitePersist) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	persist = new(SQLitePersist)

	persist.db = db
	return persist
}

func (p *SQLitePersist) InitDB() {
	if _, err := p.db.Exec(initQuery); err != nil {
		log.Fatal(err)
	}
}

func (p *SQLitePersist) getSalt(username string) (salt []byte) {
	if p.saltFetchStmt == nil {
		var err error
		p.saltFetchStmt, err = p.db.Prepare(saltFetchQuery)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := p.saltFetchStmt.QueryRow(username).Scan(&salt); err != nil {
		log.Print(err)
		return nil
	}
	return
}

func (p *SQLitePersist) Authenticate(ld data.LoginData) persistence.Session {
	salt := p.getSalt(ld.Username)

	if p.sessionCreateStmt == nil {
		var err error
		p.sessionCreateStmt, err = p.db.Prepare(sessionCreateQuery)
		if err != nil {
			log.Fatal(err)
		}
	}

	rs, err := p.sessionCreateStmt.Exec(ld.Username, salt, ld.Hash(salt))
	if err != nil {
		log.Println(err)
	}
	var rowid int64
	if rowid, err = rs.LastInsertId(); err != nil {
		log.Print(err)
		return nil
	}
	return &Session{p, uint64(rowid)}
	/*
		var session data.Session
		session.User = new(data.User)


		if err != nil {
			log.Println(err)
			return nil
		}
		tx.Commit()
		return &session

	*/

}

func (p *SQLitePersist) Register(ld data.LoginData, email string) persistence.User {
	salt := make([]byte, 36)
	_, err := rand.Read(salt)
	if err != nil {
		log.Println(err)
		return nil
	}
	hash := ld.Hash(salt)
	if err != nil {
		log.Println(err)
		return nil
	}
	if p.userRegisterStmt == nil {
		p.userRegisterStmt, err = p.db.Prepare(userRegisterQuery)
		if err != nil {
			log.Fatal(err)
		}
	}
	rs, err := p.userRegisterStmt.Exec(ld.Username, email, salt, hash)
	if err != nil {
		log.Println(err)
		return nil
	}
	rowid, err := rs.LastInsertId()
	if err != nil {
		log.Println(err)
		return nil
	}

	/*
	 */
	return &User{p, uint64(rowid)}
}

func (p *SQLitePersist) GetSession(token string) persistence.Session {
	if p.sessionIdFromTokenStmt == nil {
		var err error
		if p.sessionIdFromTokenStmt, err = p.db.Prepare(sessionIdFromTokenQuery); err != nil {
			log.Fatal(err)
		}
	}
	var sessionId uint64
	if err := p.sessionIdFromTokenStmt.QueryRow(token).Scan(&sessionId); err != nil {
		log.Println(err)
		return nil
	}
	
	s :=&Session{SQLitePersist:p, sessionid:sessionId}
	log.Printf("%#+v\n", s)
	return s
}

func (p *SQLitePersist) GetUser(username string) persistence.User {
	log.Println("seeking user",username)
	if p.userIdFromNameStmt == nil {
		var err error
		if p.userIdFromNameStmt, err = p.db.Prepare(userIdFromNameQuery); err != nil {
			log.Fatal(err)
		}
	}
	var userId uint64
	if err := p.userIdFromNameStmt.QueryRow(username).Scan(&userId); err != nil {
		log.Println(err)
		return nil
	}
	log.Printf("returning user id: %#+v\n",userId)
	return &User{p, userId}
}

func (s *Session) Get() *data.Session {
	if s.sessionFromIdStmt == nil {
		var err error
		s.sessionFromIdStmt, err = s.db.Prepare(sessionFromIdQuery)
		if err != nil {
			log.Fatal(err)
		}
	}

	sessiondata := new(data.Session)
	sessiondata.User = new(data.User)
	var lastActive int64
	var expires int64

	err := s.sessionFromIdStmt.QueryRow(s.sessionid).Scan(&(sessiondata.User.Username),
		&(sessiondata.User.Displayname), &(sessiondata.User.Email), &(sessiondata.Token),
		&lastActive, &expires)
	if err != nil {
		log.Println(err)
		return nil
	}
	sessiondata.LastActive = time.Unix(lastActive, 0)
	sessiondata.Expiry = time.Unix(expires, 0)
	return sessiondata
}

func (s *Session) User() persistence.User {
	return s.GetUser(s.Get().User.Username)
}

func (s *Session) Valid() bool {
	if s == nil {
		return false
	}
	var exp uint64
	if s.sessionExpiryStmt == nil {
		var err error
		if s.sessionExpiryStmt, err = s.db.Prepare(sessionExpiryQuery); err != nil {
			log.Fatal(err)
		}
	}

	if err := s.sessionExpiryStmt.QueryRow(s.sessionid).Scan(&exp); err != nil {
		return false
	}
	if time.Unix(int64(exp), 0).After(time.Now()) {
		return true
	}
	return false
}

func (s *Session) Invalidate() {
	if s.sessionInvalidateStmt == nil {
		var err error
		if s.sessionInvalidateStmt, err = s.db.Prepare(sessionInvalidateQuery); err != nil {
			log.Fatal(err) 
		}
	}

	if _, err := s.sessionInvalidateStmt.Exec(s.sessionid); err != nil {
		log.Println(err)
	}
}

func (u *User) Get() *data.User {
	if u.userFromIdStmt == nil {
		var err error
		u.userFromIdStmt, err = u.db.Prepare(userFromIdQuery)
		if err != nil {
			log.Fatal(err)
		}
	}
	var user data.User
	err := u.userFromIdStmt.QueryRow(u.userid).Scan(&(user.Username), &(user.Displayname), &(user.Email))
	if err != nil {
		log.Println(err)
		return nil
	}

	return &user
}

func (u *User) SetUsername(username string) persistence.User {
	if u.setUsernameStmt == nil {
		var err error
		if u.setUsernameStmt, err = u.db.Prepare(setUsernameQuery); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := u.setUsernameStmt.Exec(username, u.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) SetEmail(email string) persistence.User {
	if u.setEmailStmt == nil {
		var err error
		if u.setEmailStmt, err = u.db.Prepare(setEmailQuery); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := u.setEmailStmt.Exec(email, u.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) SetDisplayname(displayname string) persistence.User {
	if u.setDisplaynameStmt == nil {
		var err error
		if u.setDisplaynameStmt, err = u.db.Prepare(setDisplaynameQuery); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := u.setDisplaynameStmt.Exec(displayname, u.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) SetPassword(password string) persistence.User {
	if u.setHashStmt == nil {
		var err error
		if u.setHashStmt, err = u.db.Prepare(setHashQuery); err != nil {
			log.Fatal(err)
		}
	}
	username := u.Get().Username
	salt := u.getSalt(username)
	ld := &data.LoginData{username, password}

	if _, err := u.setHashStmt.Exec(ld.Hash(salt), u.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) AddFriend(friend persistence.User) persistence.User {
	f := friend.(*User)
	if u.addFriendStmt == nil {
		var err error
		if u.addFriendStmt, err = u.db.Prepare(addFriendQuery); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := u.addFriendStmt.Exec(u.userid, f.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) ConfirmFriend(friend persistence.User) persistence.User {
	f := friend.(*User)
	if u.confirmFriendStmt == nil {
		var err error
		if u.addFriendStmt, err = u.db.Prepare(confirmFriendQuery); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := u.confirmFriendStmt.Exec(u.userid, f.userid); err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func (u *User) Friends() []persistence.User {
	if u.getFriendsStmt == nil {
		var err error
		if u.getFriendsStmt, err = u.db.Prepare(getFriendsQuery); err != nil {
			log.Fatal(err)
		}
	}
	friends := make([]persistence.User, 0)
	var rows *sql.Rows
	var err error
	if rows, err = u.getFriendsStmt.Query(u.userid); err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		friend := new(User)
		friend.SQLitePersist = u.SQLitePersist
		if err := rows.Scan(&((*friend).userid)); err != nil {
			log.Println(err)
			return nil
		}
		friends = append(friends, friend)
	}
	return friends
}

func (u *User) Log(substance, route string, dose int, unit string, time time.Time) {
	if u.logDoseStmt == nil {
		var err error
		if u.logDoseStmt, err = u.db.Prepare(logDoseQuery); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := u.logDoseStmt.Exec(u.userid, substance, route, dose, unit, time.Unix()); err != nil {
		log.Println(err)
		return 
	}
}

func (u *User) History() persistence.Doses {
	doses := new(Doses)
	doses.SQLitePersist = u.SQLitePersist
	doses.query = getHistoryQuery
	doses.args = []interface{}{u.userid}
	return doses
}

func (d *Doses) Before(t time.Time) persistence.Doses {
	return &Doses{SQLitePersist: d.SQLitePersist,
		query: d.query + " AND taken < ?",
		args:  append(append([]interface{}{}, d.args...), t.Unix())}
}

func (d *Doses) After(t time.Time) persistence.Doses {
	return &Doses{SQLitePersist: d.SQLitePersist,
		query: d.query + " AND taken > ?",
		args:  append(append([]interface{}{}, d.args...), t.Unix())}
}

func (d *Doses) OfSubstance(substance string) persistence.Doses {
	return &Doses{SQLitePersist: d.SQLitePersist,
		query: d.query + " AND substance = ?",
		args:  append(append([]interface{}{}, d.args...), substance)}
}

func (d *Doses) LastX(x uint64) persistence.Doses {
	return &Doses{SQLitePersist: d.SQLitePersist,
		query: "SELECT * FROM (" + d.query + " ORDER BY taken DESC LIMIT ?) WHERE 1=1", //so chains can be added
		args:  append(append([]interface{}{}, d.args...), x)}
}

func (d *Doses) Get() []persistence.Dose {
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(d.query, d.args...); err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	doses := make([]persistence.Dose, 0)
	for rows.Next() {
		dose := new(Dose)
		dose.SQLitePersist = d.SQLitePersist
		if err = rows.Scan(&(dose.doseid)); err != nil {
			log.Println(err)
			return nil
		}
		doses = append(doses, dose)
	}
	log.Printf("%#+v",doses)
	return doses
}

func (d *Dose) SetWhen(t time.Time) persistence.Dose {
	if d.doseSetWhenStmt == nil {
		var err error
		if d.doseSetWhenStmt, err = d.db.Prepare(doseSetWhenQuery); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := d.doseSetWhenStmt.Exec(t.Unix(), d.doseid); err != nil {
		log.Println(err)
		return nil
	}
	return d
}

func (d *Dose) SetAmount(amount int) persistence.Dose {
	if d.doseSetAmountStmt == nil {
		var err error
		if d.doseSetAmountStmt, err = d.db.Prepare(doseSetAmountQuery); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := d.doseSetAmountStmt.Exec(amount, d.doseid); err != nil {
		log.Println(err)
		return nil
	}
	return d
}

func (d *Dose) SetSubstance(substance string) persistence.Dose {
	if d.doseSetSubstanceStmt == nil {
		var err error
		if d.doseSetSubstanceStmt, err = d.db.Prepare(doseSetSubstanceQuery); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := d.doseSetSubstanceStmt.Exec(substance, d.doseid); err != nil {
		log.Println(err)
		return nil
	}
	return d
}

func (d *Dose) SetRoute(route string) persistence.Dose {
	if d.doseSetRouteStmt == nil {
		var err error
		if d.doseSetRouteStmt, err = d.db.Prepare(doseSetRouteQuery); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := d.doseSetRouteStmt.Exec(route, d.doseid); err != nil {
		log.Println(err)
		return nil
	}
	return d
}

func (d *Dose) Get() *data.Dose {
	if d.getDoseStmt == nil {
		var err error
		if d.getDoseStmt, err = d.db.Prepare(getDoseQuery); err != nil {
			log.Fatal(err)
		}
	}

	var userid uint64
	var timestamp int64
	var dose data.Dose

	if err := d.getDoseStmt.QueryRow(d.doseid).Scan(&userid, &(dose.Substance), &(dose.Route), &(dose.Amount),&(dose.Unit), &timestamp); err != nil {
		log.Println(err)
		return nil
	}
	dose.User = (&User{SQLitePersist: d.SQLitePersist, userid: userid}).Get()
	dose.Taken = time.Unix(timestamp, 0)
	log.Printf("%#+v",dose)
	return &dose
}
