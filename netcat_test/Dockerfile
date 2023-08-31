FROM ubuntu:latest
RUN apt-get update && apt-get install -y netcat
COPY ./netcat_start.sh /netcat_start.sh
ENTRYPOINT [ "bash" ]
CMD [ "/netcat_start.sh" ]