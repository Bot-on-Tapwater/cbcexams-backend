-- Populate table with dump file
psql -U capeatlas -d cbcbackend -f ~/resources.sql

-- Create SQL dump without creating table
pg_dump -U capeatlas -d cbcbackend \
  --data-only \
  --no-owner \
  --inserts \
  --on-conflict-do-nothing \
  -t web_crawler_resources \
  > resources.sql

-- Full dump with table creation
pg_dump -U capeatlas -d cbcbackend -t web_crawler_resources -f resources.sql

-- Dump with date of the dump
pg_dump -U capeatlas -d cbcbackend \
  --file="web_crawler_backup_$(date +'%A_%d_%B_%Y_at_%I-%M%p').sql" \
  -t web_crawler_resources

-- Not SQL command, it's a GoLang command for linting
gofmt -s -w .