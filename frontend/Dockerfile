FROM node:20-alpine AS deps

RUN corepack enable && corepack prepare pnpm@9.15.4 --activate

WORKDIR /app

COPY app/package.json app/pnpm-lock.yaml ./

RUN pnpm fetch --prod=false
RUN pnpm install --frozen-lockfile --ignore-scripts

FROM deps AS codegen

ARG SYNC_SWAGGER="false"

COPY app/ ./
COPY swagger/ /swagger/
COPY scripts/ /scripts/

WORKDIR /app

RUN if [ "$SYNC_SWAGGER" = "true" ]; then \
      pnpm exec tsx /scripts/sync-swagger.ts || echo "sync-swagger failed, using committed specs"; \
    fi \
 && pnpm codegen

FROM deps AS builder

ARG VITE_API_URL=""
ARG VITE_USE_MSW="false"

COPY app/ ./
COPY --from=codegen /app/src/api/generated ./src/api/generated

WORKDIR /app

RUN pnpm build

FROM nginx:alpine AS production

COPY --from=builder /app/dist /usr/share/nginx/html
# Use nginx template — /docker-entrypoint.d/20-envsubst-on-templates.sh
# substitutes ${VAR} from container env into /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Defaults for docker-compose network. Override via `environment:` to point
# any service to a different host (e.g. https://application.example.com).
ENV AUTH_URL=http://auth-service:8000 \
    CANDIDATE_URL=http://candidate-service:8081 \
    COMPANY_URL=http://company-service:8080 \
    APPLICATION_URL=http://application-service:8080 \
    NGINX_ENVSUBST_TEMPLATE_SUFFIX=.template

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
