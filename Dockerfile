FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o api ./cmd/server/


FROM ubuntu:19.04 AS release

ENV PGVER 11
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 32MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 768MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

EXPOSE 5000

WORKDIR /app
COPY --from=builder /src/api .
COPY --from=builder /src/build ./build

USER postgres
CMD service postgresql start && ./api
