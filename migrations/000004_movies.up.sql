CREATE TABLE "tickitz".movies (
	id serial4 NOT NULL,
	small_photo_movie text NULL,
	big_photo_movie text NULL,
	movie_name varchar(255) NOT NULL,
	release_date date NOT NULL,
	directed_by varchar(75) NOT NULL,
	duration interval NOT NULL,
	sinopsis text NOT NULL,
	movie_cast varchar(255) NOT NULL,
	genre varchar(100) NOT NULL,
	categories varchar(25) NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT movies_categories_check CHECK (((categories)::text = ANY (ARRAY[('G'::character varying)::text, ('PG'::character varying)::text, ('PG-13'::character varying)::text, ('R'::character varying)::text, ('NC-17'::character varying)::text]))),
	CONSTRAINT pk_movies PRIMARY KEY (id)
);