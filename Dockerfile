FROM alpine:3.18 AS workspaceBuilder
COPY .build/batchv1-controller /batchv1-controller
RUN chmod +x /batchv1-controller

FROM scratch
COPY --from=workspaceBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=workspaceBuilder /batchv1-controller /batchv1-controller

WORKDIR "/"
ENTRYPOINT ["/batchv1-controller"]
