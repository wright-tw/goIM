FROM alpine:latest

RUN apk add tzdata
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN ln -snf /usr/share/zoneinfo/Asia/Taipei /etc/localtime && echo Asia/Taipei > /etc/timezone
WORKDIR /app
COPY goIM .
COPY ./src ./src

CMD ["./goIM"]