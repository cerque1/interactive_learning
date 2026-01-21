CREATE TABLE IF NOT EXISTS public.selected_categories
(
    user_id integer NOT NULL,
    category_id integer NOT NULL,
    CONSTRAINT selected_categories_pkey PRIMARY KEY (user_id, category_id)
);

CREATE TABLE IF NOT EXISTS public.selected_modules
(
    user_id integer NOT NULL,
    module_id integer NOT NULL,
    CONSTRAINT selected_modules_pkey PRIMARY KEY (user_id, module_id)
);

ALTER TABLE IF EXISTS public.selected_categories
    ADD CONSTRAINT selected_categories_category_id_fkey FOREIGN KEY (category_id)
    REFERENCES public.categories (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.selected_categories
    ADD CONSTRAINT selected_categories_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.selected_modules
    ADD CONSTRAINT selected_modules_module_id_fkey FOREIGN KEY (module_id)
    REFERENCES public.modules (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.selected_modules
    ADD CONSTRAINT selected_modules_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;