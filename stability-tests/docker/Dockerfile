ARG Lings_IMAGE
ARG LingsMINER_IMAGE

FROM ${Lings_IMAGE} as lings
FROM ${LingsMINER_IMAGE} as lingsminer

FROM golang:1.19-alpine

RUN mkdir -p /go/src/github.com/ammm56/lings

WORKDIR /go/src/github.com/ammm56/lings

RUN apk add bash build-base git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY --from=lings /app/ /app/
COPY --from=lingsminer /app/ /app/
ENV PATH="/app:${PATH}"

WORKDIR /go/src/github.com/ammm56/lings/stability-tests

RUN go install ./...

ENTRYPOINT ["./run/run.sh"]