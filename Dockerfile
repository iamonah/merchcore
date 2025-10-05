ARG G0_VERSION=1.25-alpine3.22
ARG OS_VERSION=3.21

FROM golang:${G0_VERSION} AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomods,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

# FROM base AS dev
# RUN go install github.com/air-verse/air@latest && \ 
#     go install github.com/go-delve/delve/cmd/dlv@latest
# CMD ["air", "-c", ".air.toml"]

FROM base AS build-production
WORKDIR /app
COPY . .
RUN --mount=type=cache,id=gomods,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -x -ldflags="-s -w" -o ./bin/main cmd/main.go

FROM alpine:${OS_VERSION} AS runtime
WORKDIR /app
COPY  --from=build-production /app/bin/main /app/
USER 1001:1001
EXPOSE 8080
CMD [ "/bin/main" ]