CREATE TABLE "tickitz".order_transaction (
	id varchar(100) NOT NULL,
	user_id int4 NOT NULL,
	schedules_id int4 NOT NULL,
	seats varchar(50) NOT NULL,
	total_ticket int4 NOT NULL,
	total_purchase int4 NOT NULL,
	paid bool NOT NULL DEFAULT false,
	activate_until date NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	expiri_time varchar(50) NULL,
	payment_link varchar(255) NULL,
	status varchar(25) NOT NULL DEFAULT 'Pending'::character varying,
	qr_code text NULL,
	CONSTRAINT order_transaction_status_check CHECK (((status)::text = ANY ((ARRAY['Pending'::character varying, 'Done'::character varying, 'Cancelled'::character varying])::text[]))),
	CONSTRAINT pk_order_transaction PRIMARY KEY (id)
);

ALTER TABLE "tickitz".order_transaction ADD CONSTRAINT order_transaction_schedules_id_fkey FOREIGN KEY (schedules_id) REFERENCES "tickitz".schedules(id);
ALTER TABLE "tickitz".order_transaction ADD CONSTRAINT order_transaction_user_id_fkey FOREIGN KEY (user_id) REFERENCES "tickitz".users(id);