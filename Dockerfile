FROM golang:1.27-alpine AS build
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY go.mod ./
COPY *.go ./
COPY web ./web
RUN go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /api .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /api /api
USER 65532:65532
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/api"]
