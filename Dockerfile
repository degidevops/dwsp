FROM golang:alpine as buider

ENV APP_NAME dwsp
ENV CMD_PATH dwsp.go

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

FROM alpine

ENV APP_NAME dwsp
COPY --from=buider /$APP_NAME .
CMD ./$APP_NAME
