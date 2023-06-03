FROM node:lts-alpine
RUN npm install -g http-server
WORKDIR /
COPY /frontend .
RUN npm install
RUN npm run build
EXPOSE 8080
CMD [ "http-server", "dist" ]
