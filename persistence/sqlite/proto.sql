PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;


CREATE TABLE user(
	userid INTEGER PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	displayname TEXT NOT NULL DEFAULT username,
	email TEXT UNIQUE NOT NULL,
	salt BLOB NOT NULL,
	hash BLOB NOT NULL
);

CREATE TABLE  session(
	sessionid INTEGER PRIMARY KEY,
	token TEXT UNIQUE DEFAULT (hex(randomblob(36))),
	userid INTEGER NOT NULL,
	last_access INTEGER NOT NULL DEFAULT (strftime('%s','now')),
	expires INTEGER GENERATED ALWAYS AS (strftime('%s', last_access,  'unixepoch',  '+30 days' )) VIRTUAL,
	FOREIGN KEY(userid) REFERENCES user(userid)
);


CREATE TABLE friends_t(
	friendid INTEGER PRIMARY KEY,
	user1 INTEGER NOT NULL,
	user2 INTEGER NOT NULL,
	confirmed INTEGER NOT NULL DEFAULT(0),
	FOREIGN KEY(user1) REFERENCES user(userid) ON DELETE CASCADE,
	FOREIGN KEY(user2) REFERENCES user(userid) ON DELETE CASCADE
);

CREATE VIEW friends(friendid, user, friend, confirmed) AS
	SELECT friendid, user1, user2, confirmed FROM friends_t 
	UNION SELECT friendid, user2, user1, confirmed FROM friends_t;



CREATE UNIQUE INDEX expiry
ON session(expires);


CREATE TABLE logentry(
	entryid INTEGER PRIMARY KEY,
	taken INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	amount INTEGER,
	unit TEXT,
	user INTEGER NOT NULL,
	substance TEXT NOT NULL,
	route TEXT NOT NULL,
	FOREIGN KEY (user) REFERENCES user(userid)
);

CREATE INDEX log_user
ON logentry(user);

CREATE INDEX log_substance
ON logentry(substance);

CREATE INDEX log_route
ON logentry(route);



COMMIT;
