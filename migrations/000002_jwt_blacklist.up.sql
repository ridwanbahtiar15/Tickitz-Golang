CREATE TABLE "tickitz".jwt_blacklist (
	id serial4 NOT NULL,
	jwt_code varchar(255) NOT NULL,
	CONSTRAINT pk_jwt_blacklist PRIMARY KEY (id)
);