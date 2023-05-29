FROM golang:1.20-buster AS build
ARG BUILD_DIR
ARG SERVICE_ID
ARG PUBSUB_CRED
ENV SERVICE_ID ${SERVICE_ID}
ENV PUBSUB_CRED ${PUBSUB_CRED}
COPY . /src
RUN if echo "${SERVICE_ID}" | grep -q "^aggregator"; then echo "$PUBSUB_CRED" > /pubsub.json; fi
WORKDIR /src
RUN go mod tidy
RUN cd ${BUILD_DIR} && go build -o /src/app

FROM debian:buster
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
ARG SERVICE_ID
ENV SERVICE_ID ${SERVICE_ID}
COPY ./web /web
COPY --from=build /*.json /
COPY --from=build /src/app /app
CMD /app