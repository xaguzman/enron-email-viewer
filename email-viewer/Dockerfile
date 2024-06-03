# Step 1: Build the application
FROM node:22-alpine as builder

WORKDIR /app

ARG VITE_API_URL

# Expose the environment variables
ENV VITE_API_URL=${VITE_API_URL}
RUN echo "API URL: $VITE_API_URL"


COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

# Step 2: Set up the server environment
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
