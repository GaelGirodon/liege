FROM golang AS build
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN bash ./scripts/build.sh && mkdir -p ./dist/data

FROM scratch
COPY --from=build /app/dist /
ENV LIEGE_ROOT=/data LIEGE_PORT=3000 LIEGE_VERBOSE=false
VOLUME ["/data"]
EXPOSE 3000
WORKDIR /
CMD ["/liege"]
