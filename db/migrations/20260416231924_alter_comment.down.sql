-- Reverse the changes from the up migration
ALTER TABLE comments DROP COLUMN IF EXISTS parent_id;

-- Or to drop a constraint:
-- ALTER TABLE comments DROP CONSTRAINT constraint_name;

-- Or to restore a column type:
-- ALTER TABLE comments MODIFY COLUMN column_name OLD_TYPE;