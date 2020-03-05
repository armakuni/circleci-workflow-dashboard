FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY circleci-workflow-dashboard .
COPY templates ./templates
COPY assets ./assets
ENV GIN_MODE release
EXPOSE 8080
CMD ["./circleci-workflow-dashboard"]
