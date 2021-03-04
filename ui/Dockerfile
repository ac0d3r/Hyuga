FROM node:lts-alpine as builder

WORKDIR /opt/ui

COPY . .
RUN npm install
RUN npm run build

#------------------
FROM nginx

COPY --from=builder /opt/ui/dist/ /var/www
COPY nginx.conf /etc/nginx/conf.d/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]