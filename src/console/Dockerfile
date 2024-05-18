# Build

FROM node:16 AS build

WORKDIR /app

RUN npm cache clean --force
COPY . .
RUN npm install
RUN npm run build --prod

# Deploy

FROM nginx:latest as nginx
COPY --from=build /app/dist/console /usr/share/nginx/html
COPY /nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80