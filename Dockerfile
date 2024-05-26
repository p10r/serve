FROM golang:1.21-alpine as build

# Supercronic
RUN apk add --no-cache --update go gcc g++ curl
# Latest releases available at https://github.com/aptible/supercronic/releases
ENV SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.1.9/supercronic-linux-amd64 \
    SUPERCRONIC=supercronic-linux-amd64 \
    SUPERCRONIC_SHA1SUM=5ddf8ea26b56d4a7ff6faecdd8966610d5cb9d85

RUN curl -fsSLO "$SUPERCRONIC_URL" \
 && echo "${SUPERCRONIC_SHA1SUM}  ${SUPERCRONIC}" | sha1sum -c - \
 && chmod +x "$SUPERCRONIC" \
 && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
 && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic

####

#Go
WORKDIR /go/src/app
COPY . .
RUN go test -v
ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
RUN CGO_ENABLED=1 GOOS=linux go build -o /go/bin/app

FROM alpine:latest

COPY --from=build /usr/local/bin/supercronic-linux-amd64 /supercronic
COPY --from=build /go/src/app/crontab /crontab
COPY --from=build /go/bin/app /

CMD ["/supercronic", "crontab"]