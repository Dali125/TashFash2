# Build stage
ARG GO_VERSION=1.22.5
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy the entire project to the build context
COPY . .

# Download dependencies
RUN go mod download -x

ARG TARGETARCH

# Build the Go application
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

# Final stage
FROM alpine:latest AS final

# Install runtime dependencies
RUN apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

# Create non-privileged user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Set the working directory to /app
WORKDIR /app

# Copy the executable from the build stage
COPY --from=build /bin/server /bin/

# Copy the frontend directory from the build stage
COPY --from=build /src/frontend /app/frontend

# Ensure the appuser has ownership of the frontend directory
RUN chown -R appuser:appuser /app/frontend

# Set the user to the non-privileged user
USER appuser

# Expose the port the app runs on
EXPOSE 4000

# Start the application
ENTRYPOINT [ "/bin/server" ]
