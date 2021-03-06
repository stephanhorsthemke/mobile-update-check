# Project: mobile-update-check
FROM golang:latest AS build-env

RUN echo $GOPATH

WORKDIR /go/src/github.com/egymgmbh/mobile-update-check
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

WORKDIR /go/src/github.com/egymgmbh/mobile-update-check
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-env /go/src/github.com/egymgmbh/mobile-update-check/app /srv/
COPY --from=build-env /go/src/github.com/egymgmbh/mobile-update-check/rules.json /srv/rules.json
EXPOSE 8080
WORKDIR /srv/
CMD ["./app"]
