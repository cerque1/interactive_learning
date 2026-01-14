ALTER TABLE public.modules_res
DROP COLUMN IF EXISTS time;

ALTER TABLE public.category_res
DROP COLUMN IF EXISTS time;

ALTER TABLE public.modules_res
ADD COLUMN IF NOT EXISTS time timestamp without time zone NOT NULL;

ALTER TABLE public.category_res
ADD COLUMN IF NOT EXISTS time timestamp without time zone NOT NULL;

