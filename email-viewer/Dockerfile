# Step 1: Build the application
FROM node:22-alpine as builder

WORKDIR /app

ARG VITE_API_URL
ARG VITE_API_USER
ARG VITE_API_USER_PWD

# Expose the environment variables
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_API_USER=${VITE_API_USER}
ENV VITE_API_USER_PWD=${VITE_API_USER_PWD}

RUN echo "API URL: $VITE_API_URL"
RUN echo "API User: $VITE_API_USER"
RUN echo "API Password: $VITE_API_USER_PWD"


COPY package.json yarn.lock ./

# Install dependencies / Build app
RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

# Step 2: Set up the server environment
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
