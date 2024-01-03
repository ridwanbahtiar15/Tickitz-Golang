CREATE TABLE "tickitz".cinemas (
	id serial4 NOT NULL,
	cinema_name varchar(100) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	cinema_logo text NULL,
	CONSTRAINT pk_cinemas PRIMARY KEY (id)
);