FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN make indexer

#FROM  gcr.io/distroless/base-debian12
FROM golang:1.22
WORKDIR /app
COPY --from=build /app/bin/ilx-indexer /app/ilx-indexer
CMD ["/app/ilx-indexer"]
