FROM golang as builder
WORKDIR /src/rms-media-descovery
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rms-media-discovery -a -installsuffix cgo rms-media-discovery.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/rms-media-descovery/rms-media-discovery .
EXPOSE 8080/tcp
EXPOSE 2112/tcp
ENV RMS_DATABASE=mongodb://localhost:27017
CMD ["sh", "-c", "./rms-media-discovery -host 0.0.0.0 -db ${RMS_DATABASE}"]