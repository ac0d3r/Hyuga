FROM scratch

EXPOSE 5000
EXPOSE 52/udp

COPY ./hyuga /

CMD ["/hyuga"]