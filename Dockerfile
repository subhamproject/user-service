FROM golang:1.19.4-buster


#copy binary
COPY bin/user-service user-service

EXPOSE 9091
# Run the binary program produced by `go install`
CMD ["./user-service"]


