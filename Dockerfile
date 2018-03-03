FROM golang

WORKDIR /opt/

ADD . /opt/

RUN /opt/build.sh

ENTRYPOINT ./bin/shkaff

EXPOSE 8080
