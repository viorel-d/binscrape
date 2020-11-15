FROM node:alpine

RUN mkdir /frontend
WORKDIR /frontend

COPY . .

RUN yarn install

CMD ["yarn", "dev"]
