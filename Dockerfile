FROM golang:alpine as buider

ENV APP_NAME dwsp
ENV CMD_PATH dwsp.go

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM alpine

RUN apk --no-cache add dropbear supervisor

COPY files/supervisor/supervisord.conf /etc/supervisord.conf
COPY files/supervisor/dropbear.ini /etc/supervisor.d/dropbear.ini
COPY files/supervisor/dwsp.ini /etc/supervisor.d/dwsp.ini
COPY files/ssh/authorized_keys /root/.ssh/authorized_keys

ENV APP_NAME dwsp
COPY --from=buider /$APP_NAME /bin/.
CMD /usr/bin/supervisord -c /etc/supervisord.conf
