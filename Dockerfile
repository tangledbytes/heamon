FROM alpine
RUN apk add ca-certificates

ENTRYPOINT [ "/heamon" ]

COPY heamon /