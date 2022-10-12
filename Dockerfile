FROM golang:1.19
ARG BUILD_VERSION=N/A
ARG BUILD_DATE=N/A
ARG BUILD_COMMIT=N/A
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X github.com/gopherlearning/gophkeeper/internal/conf.buildCommit=${BUILD_COMMIT} -X github.com/gopherlearning/gophkeeper/internal/conf.buildVersion=${BUILD_VERSION} -X github.com/gopherlearning/gophkeeper/internal/conf.buildDate=${BUILD_DATE}" -trimpath -o gophkeeper cmd/main.go

FROM alpine
ARG BUILD_VERSION=N/A
ARG BUILD_DATE=N/A
ARG BUILD_COMMIT=N/A
LABEL buildCommit=$BUILD_COMMIT
LABEL buildVersion=$BUILD_VERSION
LABEL buildDate=$BUILD_DATE
WORKDIR /app
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /app/gophkeeper /app/gophkeeper
EXPOSE 6220 9100
ENTRYPOINT [ "/app/gophkeeper"]
