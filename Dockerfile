# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/atla/athena
ADD ./model /go/src/github.com/atla/athena/model
ADD ./dba /go/src/github.com/atla/athena/dba

# COPY ./public /go/src/github.com/atla/lotd/public

# Build the lotd command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

RUN go get github.com/globalsign/mgo
RUN go get github.com/gorilla/mux
RUN go install github.com/atla/athena

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/athena

# Document that the service listens on port 8080.
EXPOSE 8000
# EXPOSE 8023