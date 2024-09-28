FROM golang:1.23 AS builder

ENV POSTGRES_USER=db
ENV POSTGRES_PASSWORD=db
ENV POSTGRES_DB=db

RUN go install github.com/sqldef/sqldef/cmd/psqldef@latest
RUN chmod ugo+x /go/bin/psqldef

FROM debian:stable-slim

COPY --from=builder /go/bin/psqldef /usr/local/bin
COPY ./schema/db/ddl.sql /ddl.sql
COPY ./docker/db-migrate/entrypoint.sh /entrypoint.sh
RUN chmod ugo+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
