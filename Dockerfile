FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY circleci-workflow-dashboard .
EXPOSE 8080
CMD ["./circleci-workflow-dashboard"]
