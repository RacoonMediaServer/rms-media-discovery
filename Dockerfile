FROM golang as builder
WORKDIR /src/rms-media-descovery
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-media-discovery -a -installsuffix cgo rms-media-discovery.go

FROM mcr.microsoft.com/playwright:v1.30.0-focal
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/rms-media-descovery/rms-media-discovery .
COPY --from=builder /src/rms-media-descovery/configs/rms-media-discovery.json /etc/rms/
EXPOSE 8080/tcp
EXPOSE 2112/tcp
CMD ["./rms-media-discovery"]