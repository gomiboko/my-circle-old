FROM node:14.17.1-alpine3.13
WORKDIR /app
COPY tsconfig.json .eslintrc.js .env.development vue.config.js ./
COPY package*.json ./
RUN npm install
COPY public/ ./public/

# GitHub Actions で単体テストを実行する為に必要なファイル
COPY tests/unit/ ./tests/unit/
COPY jest.config.js ./
COPY babel.config.js ./

COPY src/ ./src/
CMD ["npm", "run", "serve"]
