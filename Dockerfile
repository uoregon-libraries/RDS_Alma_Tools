From golang:1.22.4-bullseye

RUN apt-get update
RUN apt-get install -y nano
WORKDIR /usr/local/src/rds_alma_tools

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

RUN apt-get update
RUN apt-get install --no-install-recommends -y php php-dom
