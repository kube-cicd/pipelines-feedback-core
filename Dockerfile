FROM alpine:3.18 AS builder
COPY .build/batchv1-controller /batchv1-controller
RUN chmod +x /batchv1-controller

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /batchv1-controller /batchv1-controller

WORKDIR "/"
USER 65161
ENTRYPOINT ["/batchv1-controller"]
