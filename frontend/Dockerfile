#### Rust builder ####

FROM rust:1.77.0-slim-buster as rust-builder

WORKDIR /build
RUN cargo install wasm-pack

# The following line is mandatory to have cargo install dependencies
RUN mkdir src && echo 'fn main() {}' > src/lib.rs
COPY turbex-crypt/Cargo.* ./

RUN cargo build --release

COPY turbex-crypt/ ./

RUN wasm-pack build --target web

#### Node builder ####

FROM node:20.11.1-slim as node-builder

WORKDIR /build/svelte
COPY frontend-svelte/package*.json ./

RUN npm ci

COPY --from=rust-builder /build/pkg /build/turbex-crypt/pkg
COPY frontend-svelte ./

RUN npm run build

#### Runtime ####

FROM nginx:1.25.4-alpine as production

WORKDIR /app

COPY nginx.conf.template /etc/nginx/templates/nginx.conf.template
COPY --from=node-builder /build/svelte/build /app

ENV NGINX_ENVSUBST_OUTPUT_DIR=/etc/nginx
ENV TURBEX_BACK_HOST=127.0.0.2
ENV TURBEX_BACK_PORT=8000

EXPOSE 80

