# -- multistage docker build: stage #1: build stage
FROM golang:1.19-alpine AS build

RUN mkdir -p /go/src/github.com/ammm56/lings

WORKDIR /go/src/github.com/ammm56/lings

RUN apk add --no-cache curl git openssh binutils gcc musl-dev

COPY go.mod .
COPY go.sum .


# Cache lings dependencies
RUN go mod download

COPY . .

RUN go build $FLAGS -o LINGS .

# --- multistage docker build: stage #2: runtime image
FROM alpine
WORKDIR /app

RUN apk add --no-cache ca-certificates tini
RUN mkdir -p /.lings && chown nobody:nobody /.lings && chmod 700 /.lings

COPY --from=build /go/src/github.com/ammm56/lings/LINGS /app/
COPY --from=build /go/src/github.com/ammm56/lings/infrastructure/config/sample-lings.conf /app/

USER nobody
ENTRYPOINT [ "/sbin/tini", "--" ]
