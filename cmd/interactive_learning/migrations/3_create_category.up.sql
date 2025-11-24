CREATE TABLE IF NOT EXISTS public.categories
(
    id serial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    owner_id integer NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.category_modules
(
    category_id integer NOT NULL,
    module_id integer NOT NULL,
    CONSTRAINT category_modules_pkey PRIMARY KEY (category_id, module_id)
);

ALTER TABLE IF EXISTS public.categories
    ADD CONSTRAINT categories_owner_id_fkey FOREIGN KEY (owner_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.category_modules
    ADD CONSTRAINT category_modules_category_id_fkey FOREIGN KEY (category_id)
    REFERENCES public.categories (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.category_modules
    ADD CONSTRAINT category_modules_module_id_fkey FOREIGN KEY (module_id)
    REFERENCES public.modules (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;