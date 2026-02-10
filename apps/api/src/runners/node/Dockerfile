FROM node:20-slim

WORKDIR /runner

COPY runner.js ./

RUN npm init -y \
 && npm install glob \
 && npm pkg set type="module"

CMD ["node", "runner.js"]
