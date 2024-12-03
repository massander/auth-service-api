FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./auth-api .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/auth-api .
# COPY config/ ..

EXPOSE 8080
ENTRYPOINT [ "/app/auth-api"]