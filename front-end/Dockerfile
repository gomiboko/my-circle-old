FROM node:14.17.1-alpine3.13
WORKDIR /app
COPY tsconfig.json .eslintrc.js ./
COPY package*.json ./
RUN npm install
COPY public/ ./public/

COPY src/ ./src/
CMD ["npm", "run", "serve"]
