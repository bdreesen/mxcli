# Mendix Model SDK Go - Development Container
FROM mcr.microsoft.com/devcontainers/go:dev-1.26-bookworm

RUN rm -f /etc/apt/sources.list.d/yarn.list

# Install Adoptium JDK 21 (required by MxBuild) and JDK 17 (required by ANTLR4)
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget apt-transport-https gpg \
    python3-pip \
    sqlite3 \
    postgresql-client \
    kafkacat \
    && wget -qO - https://packages.adoptium.net/artifactory/api/gpg/key/public | gpg --dearmor -o /etc/apt/keyrings/adoptium.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/adoptium.gpg] https://packages.adoptium.net/artifactory/deb bookworm main" > /etc/apt/sources.list.d/adoptium.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
       openjdk-17-jre-headless \
       temurin-21-jdk \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install ANTLR4 tools
RUN pip3 install --break-system-packages antlr4-tools

# Install bun (required for VS Code extension builds)
# Install to /usr/local so all users (including vscode) can execute it
USER vscode
RUN BUN_INSTALL=/usr/local curl -fsSL https://bun.sh/install | bash

USER root
# Create libs directory for reflection data
RUN mkdir -p /workspaces/ModelSDKGo/libs
