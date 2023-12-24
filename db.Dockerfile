FROM postgres:15.3 as dumper

COPY migrations/*.sql /docker-entrypoint-initdb.d/


ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB ail
ENV POSTGRES_PORT 5432

EXPOSE 5432