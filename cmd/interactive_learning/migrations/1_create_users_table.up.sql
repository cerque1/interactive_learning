CREATE TABLE IF NOT EXISTS public.users
(
    id serial NOT NULL,
    login character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT login_unique UNIQUE (login)
);