FROM golang:1.19 AS build 

WORKDIR /srv/build 

COPY go.sum go.mod Makefile ./

RUN make deps 

COPY . . 

RUN make build 

FROM ubuntu:20.04 AS final 

COPY --from=build /srv/build/smssender /bin/
COPY --from=build /srv/build/wait-for-it.sh /bin/
COPY --from=build /srv/build/config.yaml /
COPY --from=build /srv/build/internal/database/migrations /internal/database/migrations

CMD = ["/bin/smssender", "serve"]