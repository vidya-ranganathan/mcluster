#just containerizing the operator
FROM ubuntu

COPY mcluster /usr/local/bin

ENTRYPOINT ["mcluster"]
