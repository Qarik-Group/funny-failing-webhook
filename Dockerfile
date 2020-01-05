FROM golang:1.13 as build

ARG OPTS

WORKDIR /buildspace
COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
#ENV GOPROXY="https://proxy.golang.org"
RUN go mod download

COPY . .

RUN VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') && \
  GIT_COMMIT=$(git rev-list -1 HEAD) && \
  env ${OPTS} CGO_ENABLED=0 GOOS=linux \
  go build -o funny-failing-webhook -v ./cmd

FROM alpine AS app
COPY --from=build /buildspace/funny-failing-webhook /usr/bin/funny-failing-webhook
EXPOSE 8080

CMD ["/usr/bin/funny-failing-webhook", "-port", "8080"]
