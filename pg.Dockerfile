FROM library/postgres
COPY dbinit.sql /docker-entrypoint-initdb.d/