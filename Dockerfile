FROM golang:alpine as buider

ENV APP_NAME dwsp
ENV CMD_PATH dwsp.go

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM alpine

RUN apk --no-cache add dropbear

RUN mkdir -p /root/.ssh \
    && chmod 0700 /root/.ssh \
    && mkdir -p /etc/dropbear \
    && passwd -u root

ENV APP_NAME dwsp
COPY --from=buider /$APP_NAME /usr/sbin/.

ENTRYPOINT [ "sh", "-c", "dropbear -Rkj -p 22; bwsp" ]

