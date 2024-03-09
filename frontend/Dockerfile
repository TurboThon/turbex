FROM node:20.11.1-slim as node-builder

WORKDIR /build
COPY frontend-svelte/package*.json .

RUN npm ci

COPY frontend-svelte .

RUN npm run build

#### Runtime ####

FROM nginx:1.25.4-alpine as production

WORKDIR /app

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=node-builder /build/build /app

