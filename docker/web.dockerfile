FROM golang:1.22 as build
WORKDIR /app
RUN apt update && apt install -y nodejs npm
COPY . .
RUN make build-deps
RUN make web

FROM  gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/bin/ilx-web /app/ilx-web
EXPOSE 3000
CMD ["/app/ilx-web"]
