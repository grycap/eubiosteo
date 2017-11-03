package state

var Schema = `

CREATE TABLE job (
	id TEXT,

	input TEXT,
	output TEXT,

	driverimage TEXT,
	zip TEXT,
	
	checkfiles BOOL,
	
	PRIMARY KEY(id)
);

CREATE TABLE image (
	id TEXT,

	name TEXT,
	path TEXT,
	
	format TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE alloc (
	id TEXT,
	
	jobid TEXT,
	workflowid TEXT,
	workflowallocid TEXT,

	input TEXT,
	output TEXT,

	status INTEGER,
	scriptpath TEXT,
	
	logoutput TEXT,
	logerror TEXT,

	error TEXT,
	
	initialtime BIGINT,
	finaltime BIGINT,
	elapsedtime BIGINT,

	PRIMARY KEY(id)
);

CREATE TABLE workflow (
	id TEXT,
	variables TEXT,
	steps TEXT,
	entry TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE workflowalloc (
	id TEXT,
	workflowid TEXT,

	status INTEGER,

	input TEXT,
	output TEXT,

	initialtime BIGINT,
	finaltime BIGINT,
	elapsedtime BIGINT,
	
	PRIMARY KEY(id)
);

`
