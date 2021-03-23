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

CREATE TABLE substance(
	substanceid INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULLi,
	unit TEXT NOT NULL DEFAULT 'mg'
);

CREATE TABLE route(
	routeid INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);


CREATE TABLE logentry(
	entryid INTEGER PRIMARY KEY,
	created INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	amount INTEGER NOT NULL,
	userid INTEGER NOT NULL,
	substanceid INTEGER NOT NULL,
	routeid INTEGER NOT NULL,
	FOREIGN KEY(userid) REFERENCES user(userid),
	FOREIGN KEY(substanceid) REFERENCES substance(substanceid),
	FOREIGN KEY(routeid) REFERENCES route(routeid)
);

CREATE INDEX log_user
ON logentry(userid);

CREATE INDEX log_substance
ON logentry(substanceid);

CREATE INDEX log_route
ON logentry(routeid);



COMMIT;
