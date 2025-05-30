init-db:
	sqlite3 db.sqlite < ./sql/schema.sql

seed-db:
	sqlite3 db.sqlite < ./sql/seed-data.sql
