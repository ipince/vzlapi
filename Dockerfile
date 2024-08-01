FROM golang:1.21 AS build

WORKDIR /go/src
COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH cd api && go build -ldflags '-extldflags "-static"' -tags osusergo,netgo -o /app ./main.go

FROM gcr.io/distroless/static:latest

EXPOSE 8080/tcp

COPY --from=build /app /app
ENTRYPOINT ["/app"]