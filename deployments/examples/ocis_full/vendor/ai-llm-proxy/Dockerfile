FROM node:24.16.0-alpine3.24 AS build
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
COPY package.json tsconfig.json ./
COPY src ./src
RUN pnpm install --ignore-scripts && pnpm build

FROM node:24.16.0-alpine3.24
WORKDIR /app
COPY --from=build /app/dist ./dist
EXPOSE 3030
CMD ["node", "dist/index.js"]
