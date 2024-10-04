ARG NODE_VERSION=20.17.0

FROM node:${NODE_VERSION}-alpine

COPY . /demo

WORKDIR /demo

EXPOSE 80

ENTRYPOINT ["/demo/entrypoint.sh"]

CMD [ "yarn", "serve" ]


