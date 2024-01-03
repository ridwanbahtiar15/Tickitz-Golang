CREATE TABLE "tickitz".schedules (
	id serial4 NOT NULL,
	movie_id int4 NOT NULL,
	price_per_ticket varchar(100) NOT NULL,
	schedule_date date NOT NULL,
	schedule_time text NULL,
	cinema_id int4 NOT NULL,
	seat_booked text NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT pk_schedules PRIMARY KEY (id)
);

ALTER TABLE "tickitz".schedules ADD CONSTRAINT schedules_cinema_id_fkey FOREIGN KEY (cinema_id) REFERENCES "tickitz".cinemas(id);
ALTER TABLE "tickitz".schedules ADD CONSTRAINT schedules_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES "tickitz".movies(id);