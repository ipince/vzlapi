FROM golang:1.21 AS build

WORKDIR /go/src
COPY . .

RUN cd api && go build -o /app ./main.go
RUN ls -la /
RUN ls -la /go/src

FROM gcr.io/distroless/static:latest

EXPOSE 8080/tcp

COPY --from=build /app /app
ENTRYPOINT ["/app"]