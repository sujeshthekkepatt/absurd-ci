FROM alpine:latest

ADD executor-scripts /usr/local/bin

RUN chmod +x /usr/local/bin/init.sh

RUN chmod +x /usr/local/bin/push-log.sh

RUN chmod +x /usr/local/bin/push-log_json.sh

# RUN chmod +x /usr/local/bin/log.js


RUN apk update

RUN apk add nodejs-current npm


RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl

RUN apk add git

RUN apk add --no-cache openssh

RUN apk add jq

RUN chmod +x ./kubectl

RUN mv ./kubectl /usr/bin/kubectl