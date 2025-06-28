# syntax=docker/dockerfile:1

FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod/sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source tree
COPY . .

# Build the API server from cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bumflix ./src/api

# --- Final stage ---
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /bumflix /bumflix

ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG AWS_DEFAULT_REGION
ARG AWS_ENDPOINT_URL_S3
ARG AWS_S3_RAW_BUCKET_NAME
ARG AWS_S3_HLS_BUCKET_NAME
ARG FRONTEND_URL

# Use them in RUN or ENV
ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
    AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
    AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION \
    AWS_ENDPOINT_URL_S3=$AWS_ENDPOINT_URL_S3 \
    AWS_S3_RAW_BUCKET_NAME=$AWS_S3_RAW_BUCKET_NAME \
    AWS_S3_HLS_BUCKET_NAME=$AWS_S3_HLS_BUCKET_NAME \
    FRONTEND_URL=$FRONTEND_URL

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

ENTRYPOINT ["/bumflix"]
