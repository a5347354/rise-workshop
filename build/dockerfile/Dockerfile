FROM golang:1.20-buster AS build
COPY . /src
WORKDIR /src
RUN go mod tidy
RUN cd cmd/api && go build -o /src/app

FROM debian:buster
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=build /src/app /app
CMD /app