FROM golang:1.18 AS build
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN mkdir -p ./dist/data && \
  go build -ldflags="-s -w" -o ./dist/liege ./cmd/liege.go

FROM scratch
COPY --from=build /app/dist /
ENV LIEGE_ROOT=/data LIEGE_PORT=3000
VOLUME ["/data"]
EXPOSE 3000
WORKDIR /
CMD ["/liege"]
