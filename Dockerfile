FROM debian:latest

COPY app /usr/local/bin/app

CMD ["app"]
EXPOSE 8001
