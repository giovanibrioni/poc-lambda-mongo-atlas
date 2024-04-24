FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS base

FROM base AS builder
WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /bin/user ./user/main.go

FROM alpine AS runner
WORKDIR /app
COPY --from=builder /bin/user /app/user
EXPOSE 8080
ENTRYPOINT ["/app/user"]