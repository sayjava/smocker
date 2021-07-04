FROM golang:1.15-alpine AS build-backend
RUN apk add --no-cache make
ARG VERSION=snapshot
ARG COMMIT
WORKDIR /go/src/github.com/Thiht/smocker
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY server/ ./server/
RUN make VERSION=$VERSION COMMIT=$COMMIT RELEASE=1 build

FROM node:14-alpine AS build-frontend
WORKDIR /wd
ENV PARCEL_WORKERS 1
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY tsconfig.json ./
COPY client/ ./client/
RUN yarn build

FROM alpine
WORKDIR /opt
EXPOSE 8080 8081
COPY ssl /opt/ssl
COPY --from=build-backend /go/src/github.com/Thiht/smocker/build/* /opt/
COPY --from=build-frontend /wd/build/* /opt/
CMD ["/opt/smocker"]
