FROM mongo:latest

EXPOSE 27017 
ENV MONGO_INITDB_ROOT_USERNAME admin
ENV MONGO_INITDB_ROOT_PASSWORD admin 
ENV MONGO_INITDB_DATABASE social_app

COPY mongo-init.js /docker-entrypoint-initdb.d/