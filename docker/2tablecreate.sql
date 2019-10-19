CREATE TABLE public.events
(
    id serial NOT NULL,
    uuid text COLLATE pg_catalog."default",
    header text COLLATE pg_catalog."default",
    datetime timestamp without time zone,
    description text COLLATE pg_catalog."default",
    owner text COLLATE pg_catalog."default",
    eventduration_start timestamp without time zone,
    eventduration_stop timestamp without time zone,
    mailingduration bigint,
    CONSTRAINT events_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.events
    OWNER to "user";