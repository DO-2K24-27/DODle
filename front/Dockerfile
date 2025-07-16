FROM node:20-alpine AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy all files
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:stable-alpine AS production

# Install envsubst for environment variable substitution
RUN apk add --no-cache gettext

# Copy built app to nginx server
COPY --from=build /app/dist /usr/share/nginx/html

# Copy custom nginx config to handle Vue Router's history mode
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Create a script to replace placeholders and start nginx
RUN echo '#!/bin/sh' > /docker-entrypoint.sh && \
    echo 'find /usr/share/nginx/html -type f -name "*.js" -o -name "*.html" -o -name "*.css" | xargs sed -i "s|PLACEHOLDER_API_URL|${API_URL}|g"' >> /docker-entrypoint.sh && \
    echo 'nginx -g "daemon off;"' >> /docker-entrypoint.sh && \
    chmod +x /docker-entrypoint.sh

# Start with the custom entrypoint
CMD ["/docker-entrypoint.sh"]
