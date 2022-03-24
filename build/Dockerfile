ARG APP_NAME

FROM golang:1.18 as builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./$APP_NAME ./cmd

FROM gcr.io/distroless/static
COPY --from=builder /app/$APP_NAME /usr/bin/$APP_NAME
ENTRYPOINT ["/usr/bin/$APP_NAME"]