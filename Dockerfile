FROM public.ecr.aws/docker/library/alpine:3.19.0

RUN mkdir /app

COPY ./bin/giproxy /app/giproxy

WORKDIR /app

EXPOSE 8080

USER 65534:65534

ENTRYPOINT ["/app/giproxy"]