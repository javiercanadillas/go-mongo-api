# Builder
FROM golang:1.19 as builder

WORKDIR /code
COPY . .

ARG SKAFFOLD_GO_GCFLAGS
RUN go mod download
RUN CGO_ENABLED=0 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath -o /app main.go

# Runner - Distroless
FROM gcr.io/distroless/static-debian11
ENV GOTRACEBACK=single
# Valid modes are file, env or api
ENV SECRETS_MODE=file
ENV SECRET_ENCRYPTION=true
# Set to release or debug
ENV GIN_MODE=debug
ENV GIN_PORT=8080
COPY --from=builder /app .
CMD ["/app"]

