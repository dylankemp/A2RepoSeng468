FROM nginx:latest

RUN apt-get update \
    && apt-get install -y git

COPY conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx","-g","daemon off;"]