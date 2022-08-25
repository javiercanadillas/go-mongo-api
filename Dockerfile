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
COPY --from=builder /app .
CMD ["/app"]
