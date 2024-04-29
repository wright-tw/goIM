FROM alpine:latest

ARG ENV
RUN apk add tzdata
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app
COPY goIM .
COPY ./src ./src

CMD ["./goIM"]