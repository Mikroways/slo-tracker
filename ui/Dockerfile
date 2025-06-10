FROM node:16 as build-deps
WORKDIR /usr/src/app
ARG API_URL
ENV REACT_APP_API_URL=$API_URL
COPY package*.json yarn.lock ./
RUN npm install --global yarn --force
COPY . ./
RUN npm install react-scripts --save -legacy-peer-deps
RUN yarn build

FROM nginx:alpine
RUN mkdir -p /app
COPY --from=build-deps /usr/src/app/build /app
COPY nginx.conf /etc/nginx/
WORKDIR /app/
CMD nginx
EXPOSE 3000