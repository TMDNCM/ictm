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

CREATE UNIQUE INDEX expiry
ON session(expires);

COMMIT;
