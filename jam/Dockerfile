FROM golang:alpine

MAINTAINER peter.trevan@bluefruit.co.uk

ENV GIN_MODE=release
ENV PORT=3004

RUN apk update && apk upgrade

WORKDIR /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/

COPY ./ /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/
#COPY ./static/* /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/
#COPY ./templates/* /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/

RUN ls /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/

RUN go get /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/

RUN go build /go/src/github.com/edmund-hazledine-bluefruit/game-jam-2021/

EXPOSE $PORT

ENTRYPOINT ["./jack-fruit"]

