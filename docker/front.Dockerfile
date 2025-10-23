FROM node:lts AS base

FROM base AS deps
WORKDIR /tmp/

COPY package.json pnpm-lock.yaml /tmp/
RUN --mount=type=cache,target=/root/.local/share/pnpm/store corepack enable pnpm && \
  pnpm install --frozen-lockfile


FROM base AS build
WORKDIR /tmp/

COPY --from=deps /tmp/node_modules ./node_modules
COPY . .

RUN corepack enable pnpm && pnpm run build

FROM cgr.dev/chainguard/nginx:latest-dev

WORKDIR /usr/share/nginx/html
USER nginx

COPY --from=build --chown=nginx:nginx /tmp/dist/  .
COPY --chown=nginx:nginx ./nginx.conf /etc/nginx/conf.d/nginx.default.conf

EXPOSE 8080
ENTRYPOINT ["/bin/sh",  "-c",  "./vite-envs.sh && exec nginx -g 'daemon off;'"]

