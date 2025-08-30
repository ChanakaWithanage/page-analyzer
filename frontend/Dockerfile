# ---- Build Stage ----
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files first for caching
COPY package*.json ./

# Install dependencies (including pinned TypeScript)
RUN npm install --legacy-peer-deps

# Copy the rest of the code
COPY . .

# Verify installed TS version (debug step)
RUN npx tsc --version

# Build
RUN npm run build

# ---- Runtime Stage ----
FROM nginx:1.25-alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/dist .
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
