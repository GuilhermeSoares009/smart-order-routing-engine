FROM golang:1.22
WORKDIR /app
CMD ["sh", "-c", "echo app service idle && sleep 3600"]
