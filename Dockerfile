FROM golang:latest as builder

LABEL app="slo-tracker"
LABEL version="0.0.1"
LABEL description="slo-tracker : Track your product SLO"

WORKDIR /app

# improves caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/slo-tracker .

FROM alpine AS final

RUN apk upgrade && \
    apk --no-cache add curl

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install atlas for db migrations
RUN curl -sSf https://atlasgo.sh | sh -s -- -y

COPY --from=builder --chown=appuser:appuser /app/slo-tracker .
COPY --chown=appuser:appuser atlas.hcl schema.hcl .
COPY --chown=appuser:appuser migrations /app/migrations

USER appuser

EXPOSE 8080
CMD /app/slo-tracker