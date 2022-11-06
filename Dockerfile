FROM golang:alpine as buider

ENV APP_NAME dwsp
ENV CMD_PATH dwsp.go

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM alpine

RUN apk --no-cache add openrc dropbear dropbear-openrc dropbear-ssh \
    && passwd -u root \
    && touch /run/openrc/softlevel \
    && rc-update add dropbear \
    && rc-service dropbear start

ENV APP_NAME dwsp
COPY --from=buider /$APP_NAME .
CMD ./$APP_NAME
