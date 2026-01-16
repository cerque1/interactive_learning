ALTER TABLE public.categories
ADD COLUMN IF NOT EXISTS type int;

UPDATE categories
SET type = 0
WHERE type IS NULL;

ALTER TABLE categories
ALTER COLUMN type SET NOT NULL;