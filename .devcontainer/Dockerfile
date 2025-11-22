# FROM mcr.microsoft.com/vscode/devcontainers/python:3.13-bookworm
FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

ARG BUF_VERSION=1.50.0
ARG HADOLINT_VERSION=2.12.0

ENV GOPRIVATE=buf.build/gen/go


SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get update \
    && apt-get install -y nodejs

# Install Buf
RUN curl -o /usr/local/bin/buf -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" \
    && chmod +x /usr/local/bin/buf

# Install Hadolint
RUN curl -o /usr/local/bin/hadolint -sSL "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-$(uname -s)-$(uname -m | sed 's/aarch64/arm64/')" \
    && chmod +x /usr/local/bin/hadolint

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
