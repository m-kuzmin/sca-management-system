ARG go_version=1.22
ARG base_image=alpine:latest


FROM sqlc/sqlc:latest as sqlc
WORKDIR /sqlc

COPY db ./db
COPY sqlc.yaml .
RUN ["/workspace/sqlc", "generate"]


FROM golang:$go_version AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --from=sqlc /sqlc/db ./db
COPY cmd ./cmd
COPY api ./api
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/rest_api/


FROM $base_image
EXPOSE 8000

COPY --from=build /app/server /usr/bin/server
COPY config/ /etc/server/
COPY db/migrations/ /var/server/migrations
CMD ["server"]
