ALTER TABLE public.modules_res
ADD COLUMN IF NOT EXISTS owner integer NOT NULL;

ALTER TABLE public.category_res
ADD COLUMN IF NOT EXISTS owner integer NOT NULL;

ALTER TABLE public.modules_res
ADD CONSTRAINT modules_res_owner_fkey 
FOREIGN KEY (owner)
REFERENCES public.users (id) MATCH SIMPLE
ON UPDATE NO ACTION
ON DELETE NO ACTION;

ALTER TABLE public.category_res
ADD CONSTRAINT category_res_owner_fkey
FOREIGN KEY (owner)
REFERENCES public.users (id) MATCH SIMPLE
ON UPDATE NO ACTION
ON DELETE NO ACTION;

ALTER TABLE IF EXISTS public.results
DROP CONSTRAINT IF EXISTS results_owner_fkey;

ALTER TABLE public.results
DROP COLUMN IF EXISTS owner;