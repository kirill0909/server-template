FROM golang:1.20.4-alpine AS build
ARG ACCESS_AX_USER
ARG ACCESS_AX_TOKEN
RUN apk --no-cache add gcc g++ make git \
    && go env -w GOPRIVATE=${CI_SERVER_HOST} GOSUMDB=off \
    && echo -e "machine gitlab.axarea.ru\nlogin ${ACCESS_AX_USER}\npassword ${ACCESS_AX_TOKEN}" > ~/.netrc

# create temp dir for build app
WORKDIR /app
ADD ./app /app

# download requirements && compile
RUN go mod download \
    && GOOS=linux GOARCH=amd64 go build -o /bin/main ./cmd/api/main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY --from=build /bin/main /bin/main
ENTRYPOINT ["/bin/main"]
