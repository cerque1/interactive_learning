CREATE TABLE IF NOT EXISTS public.cards_results
(
    result_id integer NOT NULL,
    card_id integer NOT NULL,
    result character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT cards_results_pkey PRIMARY KEY (result_id, card_id)
);

CREATE TABLE IF NOT EXISTS public.category_res
(
    category_result_id integer NOT NULL,
    category_id integer NOT NULL,
    module_id integer NOT NULL,
    result_id integer NOT NULL,
    CONSTRAINT category_res_pkey PRIMARY KEY (category_result_id, module_id)
);

CREATE TABLE IF NOT EXISTS public.modules_res
(
    module_id integer NOT NULL,
    result_id integer NOT NULL,
    CONSTRAINT modules_res_pkey PRIMARY KEY (module_id, result_id)
);

CREATE TABLE IF NOT EXISTS public.results
(
    id serial NOT NULL,
    owner integer NOT NULL,
    type character varying COLLATE pg_catalog."default" NOT NULL,
    "time" timestamp without time zone NOT NULL,
    CONSTRAINT results_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.cards_results
    ADD CONSTRAINT cards_results_card_id_fkey FOREIGN KEY (card_id)
    REFERENCES public.cards (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE IF EXISTS public.cards_results
    ADD CONSTRAINT cards_results_result_id_fkey FOREIGN KEY (result_id)
    REFERENCES public.results (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE IF EXISTS public.category_res
    ADD CONSTRAINT category_res_category_id_fkey1 FOREIGN KEY (category_id)
    REFERENCES public.categories (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.category_res
    ADD CONSTRAINT category_res_module_id_fkey1 FOREIGN KEY (module_id)
    REFERENCES public.modules (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.category_res
    ADD CONSTRAINT category_res_result_id_fkey FOREIGN KEY (result_id)
    REFERENCES public.results (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE IF EXISTS public.modules_res
    ADD CONSTRAINT modules_res_module_id_fkey FOREIGN KEY (module_id)
    REFERENCES public.modules (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.modules_res
    ADD CONSTRAINT modules_res_result_id_fkey FOREIGN KEY (result_id)
    REFERENCES public.results (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.results
    ADD CONSTRAINT results_owner_fkey FOREIGN KEY (owner)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;