# Build

FROM golang:alpine AS build
ADD . /app
RUN cd /app && go build -o server

# Deploy

FROM alpine
WORKDIR /server
COPY --from=build /app /server

EXPOSE 8082
ENTRYPOINT ["./server"]
