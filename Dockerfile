FROM golang:1.23-alpine3.20 AS base
ENV AZ_INSTALLER=DOCKER
RUN apk add py3-pip && \
    apk add --virtual=build gcc musl-dev python3-dev libffi-dev openssl-dev cargo make && \
    python3 -m venv /azure-cli-venv && \
    . /azure-cli-venv/bin/activate && \
    pip install --no-cache-dir azure-cli && \
    deactivate && \
    apk del --purge build
WORKDIR /app
ENV AZURE_CONFIG_DIR=/app/.azure

FROM golang:1.23-alpine3.20 AS publish
COPY src /src
WORKDIR /src
RUN mkdir -p /app/publish
RUN go mod download
ARG VERSION=1
RUN ls -la /src
RUN go build main.go

FROM base AS azcliproxy-dist
WORKDIR /app
COPY --from=publish /app/publish .

ENV PATH="/azure-cli-venv/bin:$PATH"
ENTRYPOINT ["./azcliproxy"]
