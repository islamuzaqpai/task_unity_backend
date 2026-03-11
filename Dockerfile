FROM postgres:16

ENV POSTGRES_USER=gouser
ENV POSTGRES_PASSWORD=gopassword
ENV POSTGRES_DB=task_unity_db

COPY init.sql /docker-entrypoint-initdb.d/