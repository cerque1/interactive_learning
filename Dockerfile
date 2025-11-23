FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/interactive_learning ./cmd/interactive_learning

FROM alpine 

WORKDIR /app

COPY /static /app/static
COPY --from=builder /app/interactive_learning /app/interactive_learning

CMD [ "./interactive_learning" ]