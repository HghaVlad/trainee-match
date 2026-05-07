FROM node:20-alpine AS deps

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

COPY app/package.json app/pnpm-lock.yaml ./

RUN pnpm fetch --prod=false
RUN pnpm install --frozen-lockfile --ignore-scripts

FROM deps AS codegen

COPY app/ ./
COPY swagger/ /swagger/
COPY scripts/ /scripts/

WORKDIR /app

RUN pnpm codegen

FROM deps AS builder

ARG VITE_API_URL=""
ARG VITE_USE_MSW="false"

COPY app/ ./
COPY --from=codegen /app/src/api/generated ./src/api/generated

WORKDIR /app

RUN pnpm build

FROM nginx:alpine AS production

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
