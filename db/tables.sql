-- This script was generated by the ERD tool in pgAdmin 4.
-- Please log an issue at https://github.com/pgadmin-org/pgadmin4/issues/new/choose if you find any bugs, including reproduction steps.
BEGIN;


CREATE TABLE IF NOT EXISTS public."Collection"
(
    id serial NOT NULL,
    collection_name character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT "Collection_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Folders"
(
    id serial NOT NULL,
    folder_name character varying(50) COLLATE pg_catalog."default",
    collection_id integer,
    CONSTRAINT "Folders_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Promotion"
(
    id serial NOT NULL,
    price character varying(255) COLLATE pg_catalog."default",
    description character varying(255) COLLATE pg_catalog."default",
    image_url character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT "Promotion_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."USER"
(
    id serial NOT NULL,
    email character varying(70) COLLATE pg_catalog."default",
    password character varying(70) COLLATE pg_catalog."default",
    user_role integer,
    CONSTRAINT "USER_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.applications
(
    id serial NOT NULL,
    date date NOT NULL,
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    phone_number character varying(20) COLLATE pg_catalog."default",
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    status character varying(20) COLLATE pg_catalog."default",
    description character varying(500) COLLATE pg_catalog."default",
    CONSTRAINT applications_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.images
(
    id serial NOT NULL,
    image_link character varying(70) COLLATE pg_catalog."default",
    collection_id integer,
    CONSTRAINT images_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.schedules
(
    id serial NOT NULL,
    date date NOT NULL,
    "time" time without time zone NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    email character varying(20) COLLATE pg_catalog."default",
    firstname character varying(20) COLLATE pg_catalog."default",
    lastname character varying(20) COLLATE pg_catalog."default",
    phone integer,
    CONSTRAINT schedules_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public."Folders"
    ADD CONSTRAINT "Folders_collection_id_fkey" FOREIGN KEY (collection_id)
    REFERENCES public."Collection" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.images
    ADD CONSTRAINT images_collection_id_fkey FOREIGN KEY (collection_id)
    REFERENCES public."Collection" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

END;