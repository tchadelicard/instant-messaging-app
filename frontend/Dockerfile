# Stage 1: Build the frontend
FROM node:20 AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Build the frontend
COPY . .
RUN npm run build

# Stage 2: Serve with Caddy
FROM caddy:2.8

# Copy the built frontend from the builder stage to Caddy's default web root
COPY --from=builder /app/dist /srv

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Expose the default Caddy port
EXPOSE 80
