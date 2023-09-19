FROM golang:1.20-alpine as builder

WORKDIR /app

COPY ./services ./services
COPY ./libs ./libs

RUN apk add rsync
RUN rsync -arm --include='*/' --include='schema.graphql' --exclude='*' . schemas

RUN go work init
RUN go work use -r .

RUN cd ./services/default/graphql && \
  CGO_ENABLED=0 \
  go build \
  -ldflags="-s -w" \
  -o /app/bin/server \
  cmd/*.go

FROM alpine:latest

RUN mkdir /schemas
ENV SCHEMA_ROOT=/schemas/

COPY --from=builder /app/schemas /schemas
COPY --from=builder /app/bin/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
