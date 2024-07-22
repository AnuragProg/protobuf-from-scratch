
ARG GO_VERSION=1.22.3

############################################
FROM golang:${GO_vERSION}-alpine AS base

WORKDIR /usr/app/src
RUN apk update & apk add make
COPY go.mod go.sum ./
RUN go mod download & go mod verify
COPY . .


############################################
FROM base AS client-build
RUN go build -o /usr/app/bin/client /usr/app/src/cmd/client/main.go


############################################
FROM base AS server-build
RUN go build -o /usr/app/bin/server /usr/app/src/cmd/server/main.go


############################################
FROM scratch AS client
WORKDIR /usr/app
COPY --from=client-build /usr/app/bin/client /bin/client
ENTRYPOINT ["/bin/client"]

############################################
FROM scratch AS server
WORKDIR /usr/app
COPY --from=server-build /usr/app/bin/server /bin/server
ENTRYPOINT ["/bin/server"]
