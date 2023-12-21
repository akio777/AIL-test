FROM postgres:15 as dumper

COPY init-db.sh /docker-entrypoint-initdb.d/
COPY migrations/*.up.sql /docker-entrypoint-initdb.d/
COPY cmd/api/.env .

EXPOSE 5432