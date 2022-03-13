FROM ubuntu
RUN apt update && apt install -y ca-certificates
COPY ./k8sgw /k8sgw
ENTRYPOINT [ "/k8sgw" ]