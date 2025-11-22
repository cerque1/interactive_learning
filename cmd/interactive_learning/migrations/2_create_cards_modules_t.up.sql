CREATE TABLE IF NOT EXISTS public.modules
(
    id serial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    owner_id integer NOT NULL,
    type character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT modules_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.modules
    ADD CONSTRAINT modules_owner_id_fkey FOREIGN KEY (owner_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

CREATE TABLE IF NOT EXISTS public.cards
(
    id serial NOT NULL,
    module_id integer NOT NULL,
    term_lang character varying COLLATE pg_catalog."default" NOT NULL,
    term_text character varying COLLATE pg_catalog."default" NOT NULL,
    def_lang character varying COLLATE pg_catalog."default" NOT NULL,
    def_text character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT cards_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.cards
    ADD CONSTRAINT module_fk FOREIGN KEY (module_id)
    REFERENCES public.modules (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

