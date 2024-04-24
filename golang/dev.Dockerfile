FROM golang:1.20-alpine AS base

FROM base AS dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
