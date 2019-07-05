FROM scratch

WORKDIR /app

ADD views /app/views

ADD gohost /

CMD ["/gohost"]
