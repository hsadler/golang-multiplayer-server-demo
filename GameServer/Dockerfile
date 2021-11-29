FROM golang:1.17.3-alpine3.14

# run the image as a non-root user
RUN adduser -D gameserveruser
USER gameserveruser

# copy over server files
COPY . /go/src/

# set working directory
WORKDIR /go/src/

# build the server executable
RUN go build

# start the server
CMD ./server