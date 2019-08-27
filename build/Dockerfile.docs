FROM node:8.11.4

WORKDIR /app/website

EXPOSE 3000
COPY ./docs /app/docs
COPY ./website /app/website
RUN yarn install

CMD ["yarn", "start", "--no-watch"]
