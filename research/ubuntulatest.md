# Ubuntu Actions Runner Image Analysis

**Last Updated**: 2026-02-07  
**Source**: [GitHub Actions Runner Images Repository](https://github.com/actions/runner-images)  
**Ubuntu Version**: 24.04 LTS (Noble Numbat)  
**Runner Label**: `ubuntu-latest`  
**Architecture**: x86_64

## Overview

This document provides an analysis of the default GitHub Actions Ubuntu runner image (`ubuntu-latest`) and guidance for creating Docker images that mimic its environment. The `ubuntu-latest` label currently points to Ubuntu 24.04 LTS.

The runner image is a comprehensive development environment with pre-installed tools, runtimes, and services commonly used in CI/CD workflows.

## Operating System

- **Distribution**: Ubuntu 24.04 LTS (Noble Numbat)
- **Kernel**: Linux 6.x (Azure-optimized)
- **Architecture**: x86_64
- **Init System**: systemd
- **Package Manager**: apt (APT 2.7.x), snap

## Language Runtimes

### Node.js
- **Versions**: Multiple versions via nvm (16.x, 18.x, 20.x, 22.x)
- **Default Version**: 20.x LTS
- **Package Managers**: 
  - npm (included with Node.js)
  - yarn 1.22.x
  - pnpm 8.x/9.x
- **Global Tools**: npx, corepack

### Python
- **Versions**: 3.10.x, 3.11.x, 3.12.x
- **Default Version**: 3.12.x
- **Package Manager**: pip 24.x
- **Additional Tools**: 
  - pipenv
  - poetry
  - virtualenv
  - pipx

### Ruby
- **Versions**: 3.0.x, 3.1.x, 3.2.x, 3.3.x
- **Default Version**: 3.2.x
- **Package Manager**: gem, bundler 2.x

### Go
- **Versions**: 1.21.x, 1.22.x, 1.23.x
- **Default Version**: 1.22.x or 1.23.x
- **Package Manager**: go modules (built-in)

### Java/JDK
- **Versions**: 
  - Temurin-11 (LTS)
  - Temurin-17 (LTS)
  - Temurin-21 (LTS)
- **Default Version**: 21
- **Build Tools**: Maven 3.9.x, Gradle 8.x, Ant 1.10.x

### PHP
- **Versions**: 8.1.x, 8.2.x, 8.3.x
- **Default Version**: 8.3.x
- **Package Manager**: composer 2.x

### Rust
- **Version**: 1.7x.x (stable channel)
- **Package Manager**: cargo (included)
- **Additional Tools**: rustup, clippy, rustfmt

### .NET
- **SDK Versions**: 6.0.x, 7.0.x, 8.0.x
- **Default Version**: 8.0.x
- **Package Manager**: dotnet (built-in)

## Container Tools

### Docker
- **Version**: 24.x or 25.x
- **Components**:
  - docker-compose v2.x
  - docker-buildx v0.12.x or later
  - docker-ce, docker-ce-cli
- **containerd**: 1.7.x
- **Storage Driver**: overlay2
- **Cgroup Driver**: systemd

### Kubernetes Tools
- **kubectl**: Latest stable (1.28+)
- **helm**: 3.14.x or later
- **minikube**: 1.32.x or later
- **kind**: 0.20.x or later

### Container Registries
- Pre-configured for Docker Hub
- Azure Container Registry (ACR) CLI support
- Amazon ECR CLI support
- Google Container Registry (GCR) support

## Build Tools

### Compilers & Build Essentials
- **gcc/g++**: 13.x
- **clang/clang++**: 18.x
- **Make**: 4.3
- **CMake**: 3.28.x or later
- **Ninja**: 1.11.x
- **Meson**: 1.3.x
- **Autoconf**: 2.71
- **Automake**: 1.16

### Additional Build Tools
- **pkg-config**: 1.8.x
- **m4**: 1.4.x
- **bison**: 3.8.x
- **flex**: 2.6.x
- **swig**: 4.1.x

## Databases & Services

### PostgreSQL
- **Version**: 14.x, 15.x, 16.x
- **Service Status**: Installed but not running by default
- **Default Port**: 5432
- **Client Tools**: psql, pg_dump, pg_restore

### MySQL
- **Version**: 8.0.x
- **Service Status**: Installed but not running by default
- **Default Port**: 3306
- **Client Tools**: mysql, mysqldump

### MongoDB
- **Version**: 7.x
- **Service Status**: Installed but not running by default
- **Default Port**: 27017
- **Client Tools**: mongosh (MongoDB Shell)

### Redis
- **Version**: 7.x
- **Service Status**: Installed but not running by default
- **Default Port**: 6379
- **Client Tools**: redis-cli

### SQLite
- **Version**: 3.45.x or later
- **Client Tools**: sqlite3

## CI/CD Tools

### GitHub CLI
- **Version**: 2.44.x or later
- **Extensions**: Pre-configured for GitHub Actions
- **Authentication**: Supports GITHUB_TOKEN by default

### Cloud CLIs

#### Azure CLI
- **Version**: 2.57.x or later
- **Extensions**: Various Azure service extensions pre-installed

#### AWS CLI
- **Version**: 2.15.x or later
- **Tools**: aws, aws-sam-cli

#### Google Cloud SDK
- **Version**: 465.x or later
- **Tools**: gcloud, gsutil, bq

### Infrastructure as Code
- **Terraform**: 1.7.x or later
- **Ansible**: 2.16.x (via pip)
- **Pulumi**: 3.x

### Other DevOps Tools
- **Jenkins**: Not pre-installed (use via Docker)
- **GitLab Runner**: Not pre-installed
- **CircleCI CLI**: Not pre-installed

## Testing Tools

### Browser Testing
- **Selenium**: 4.x (via pip/npm)
- **Playwright**: Latest (via npm)
- **Cypress**: Latest (via npm)
- **ChromeDriver**: Latest stable
- **GeckoDriver**: Latest stable

### Browsers
- **Chrome/Chromium**: Latest stable
- **Firefox**: Latest stable
- **Microsoft Edge**: Latest stable (Chromium-based)

### Unit Testing Frameworks
- **Jest**: Latest (via npm)
- **Mocha**: Latest (via npm)
- **pytest**: Latest (via pip)
- **JUnit**: Via Maven/Gradle
- **RSpec**: Via gem

## Version Control

### Git
- **Version**: 2.43.x or later
- **LFS**: git-lfs 3.4.x
- **Tools**: git-flow, git-cola (GUI)

### Other VCS
- **Subversion (SVN)**: 1.14.x
- **Mercurial**: 6.x

## Package Managers

- **apt**: 2.7.x (Debian package manager)
- **snap**: 2.61.x (Snap package manager)
- **pip**: 24.x (Python packages)
- **npm**: 10.x (Node.js packages)
- **yarn**: 1.22.x (Node.js packages)
- **pnpm**: 8.x/9.x (Node.js packages)
- **gem**: 3.5.x (Ruby packages)
- **cargo**: 1.7x.x (Rust packages)
- **composer**: 2.x (PHP packages)

## Environment Variables

Key environment variables set in the runner:

``````bash
# System paths
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin

# GitHub Actions context
GITHUB_WORKSPACE=/home/runner/work/{repo}/{repo}
GITHUB_ACTION=__run
GITHUB_ACTOR={username}
GITHUB_REPOSITORY={owner}/{repo}
GITHUB_EVENT_NAME={event}
GITHUB_SHA={commit-sha}
GITHUB_REF={ref}
GITHUB_HEAD_REF={branch}  # For PRs
GITHUB_BASE_REF={branch}  # For PRs
RUNNER_OS=Linux
RUNNER_ARCH=X64
RUNNER_NAME={runner-name}
RUNNER_TEMP=/home/runner/work/_temp

# CI flags
CI=true
GITHUB_ACTIONS=true

# Tool paths
JAVA_HOME=/usr/lib/jvm/temurin-21-jdk-amd64
ANDROID_HOME=/usr/local/lib/android/sdk
GOROOT=/opt/hostedtoolcache/go/{version}/x64
``````

## Notable Pre-installed Tools

### Code Analysis & Linting
- **ESLint**: Latest (via npm)
- **Pylint**: Latest (via pip)
- **RuboCop**: Latest (via gem)
- **ShellCheck**: 0.9.x
- **hadolint**: Latest (Dockerfile linter)

### Documentation
- **Doxygen**: 1.9.x
- **Sphinx**: Latest (via pip)
- **Jekyll**: Latest (via gem)

### Utilities
- **curl**: 8.5.x or later
- **wget**: 1.21.x
- **jq**: 1.7.x
- **yq**: 4.x (YAML processor)
- **zip/unzip**: 6.x/6.x
- **7zip**: 23.x
- **rsync**: 3.2.x
- **ssh**: OpenSSH 9.x

## Creating a Docker Image Mimic

To create a Docker image that mimics the GitHub Actions Ubuntu runner environment:

### Base Image

Start with the Ubuntu base image matching the runner version:

``````dockerfile
FROM ubuntu:24.04

# Set environment to non-interactive to avoid prompts
ENV DEBIAN_FRONTEND=noninteractive
``````

### System Setup

``````dockerfile
# Update system packages
RUN apt-get update && apt-get upgrade -y

# Install build essentials
RUN apt-get install -y \
    build-essential \
    cmake \
    git \
    git-lfs \
    curl \
    wget \
    jq \
    zip \
    unzip \
    tar \
    gzip \
    bzip2 \
    xz-utils \
    ca-certificates \
    gnupg \
    lsb-release \
    software-properties-common \
    apt-transport-https

# Install additional utilities
RUN apt-get install -y \
    vim \
    nano \
    rsync \
    openssh-client \
    p7zip-full \
    shellcheck
``````

### Language Runtimes

``````dockerfile
# Install Node.js (using NodeSource)
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

# Install Yarn and pnpm
RUN npm install -g yarn pnpm

# Install Python
RUN apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \
    python-is-python3

# Install Python tools
RUN pip3 install --no-cache-dir \
    pipenv \
    poetry \
    virtualenv \
    pytest \
    pylint

# Install Ruby
RUN apt-get install -y ruby-full && \
    gem install bundler

# Install Go
RUN wget -q https://go.dev/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz && \
    rm go1.23.1.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt

# Install Java (Temurin)
RUN wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public | apt-key add - && \
    echo "deb https://packages.adoptium.net/artifactory/deb $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/adoptium.list && \
    apt-get update && \
    apt-get install -y temurin-21-jdk

# Set JAVA_HOME
ENV JAVA_HOME=/usr/lib/jvm/temurin-21-jdk-amd64
ENV PATH="${JAVA_HOME}/bin:${PATH}"

# Install Maven
RUN wget -q https://dlcdn.apache.org/maven/maven-3/3.9.6/binaries/apache-maven-3.9.6-bin.tar.gz && \
    tar -C /opt -xzf apache-maven-3.9.6-bin.tar.gz && \
    ln -s /opt/apache-maven-3.9.6/bin/mvn /usr/local/bin/mvn && \
    rm apache-maven-3.9.6-bin.tar.gz

# Install Gradle
RUN wget -q https://services.gradle.org/distributions/gradle-8.6-bin.zip && \
    unzip -q gradle-8.6-bin.zip -d /opt && \
    ln -s /opt/gradle-8.6/bin/gradle /usr/local/bin/gradle && \
    rm gradle-8.6-bin.zip

# Install PHP
RUN apt-get install -y \
    php \
    php-cli \
    php-mbstring \
    php-xml \
    php-curl \
    php-zip && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Install .NET SDK
RUN wget -q https://packages.microsoft.com/config/ubuntu/24.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb && \
    dpkg -i packages-microsoft-prod.deb && \
    rm packages-microsoft-prod.deb && \
    apt-get update && \
    apt-get install -y dotnet-sdk-8.0
``````

### Container Tools

``````dockerfile
# Install Docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl

# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
``````

### Additional Tools

``````dockerfile
# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list && \
    apt-get update && \
    apt-get install -y gh

# Install Azure CLI
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Install AWS CLI
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf aws awscliv2.zip

# Install Google Cloud SDK
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update && \
    apt-get install -y google-cloud-sdk

# Install Terraform
RUN wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list && \
    apt-get update && \
    apt-get install -y terraform
``````

### Databases

``````dockerfile
# Install PostgreSQL client
RUN apt-get install -y postgresql-client

# Install MySQL client
RUN apt-get install -y mysql-client

# Install MongoDB client
RUN wget -qO - https://www.mongodb.org/static/pgp/server-7.0.asc | apt-key add - && \
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/7.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-7.0.list && \
    apt-get update && \
    apt-get install -y mongodb-mongosh

# Install Redis client
RUN apt-get install -y redis-tools

# Install SQLite
RUN apt-get install -y sqlite3
``````

### Environment Configuration

``````dockerfile
# Set environment variables to match runner
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin
ENV CI=true
ENV DEBIAN_FRONTEND=noninteractive

# Create a user similar to the runner user
RUN useradd -m -s /bin/bash runner && \
    mkdir -p /home/runner/work && \
    chown -R runner:runner /home/runner

# Switch to runner user
USER runner
WORKDIR /home/runner

# Set default shell
SHELL ["/bin/bash", "-c"]
``````

### Complete Dockerfile Example

Here's a complete, working Dockerfile that can be used as a starting point:

``````dockerfile
FROM ubuntu:24.04

# Prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Update and install base packages
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y \
    build-essential \
    cmake \
    git \
    git-lfs \
    curl \
    wget \
    jq \
    zip \
    unzip \
    tar \
    ca-certificates \
    gnupg \
    lsb-release \
    software-properties-common \
    apt-transport-https \
    vim \
    rsync \
    openssh-client \
    p7zip-full \
    shellcheck

# Install Node.js 20
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g yarn pnpm

# Install Python
RUN apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \
    python-is-python3 && \
    pip3 install --no-cache-dir pipenv poetry virtualenv pytest

# Install Ruby
RUN apt-get install -y ruby-full && \
    gem install bundler

# Install Go
RUN wget -q https://go.dev/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz && \
    rm go1.23.1.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Java (Temurin 21)
RUN wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public | apt-key add - && \
    echo "deb https://packages.adoptium.net/artifactory/deb $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/adoptium.list && \
    apt-get update && \
    apt-get install -y temurin-21-jdk
ENV JAVA_HOME=/usr/lib/jvm/temurin-21-jdk-amd64
ENV PATH="${JAVA_HOME}/bin:${PATH}"

# Install Maven
RUN wget -q https://dlcdn.apache.org/maven/maven-3/3.9.6/binaries/apache-maven-3.9.6-bin.tar.gz && \
    tar -C /opt -xzf apache-maven-3.9.6-bin.tar.gz && \
    ln -s /opt/apache-maven-3.9.6/bin/mvn /usr/local/bin/mvn && \
    rm apache-maven-3.9.6-bin.tar.gz

# Install PHP and Composer
RUN apt-get install -y php php-cli php-mbstring php-xml php-curl php-zip && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Install Docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Install kubectl and Helm
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl && \
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list && \
    apt-get update && \
    apt-get install -y gh

# Install database clients
RUN apt-get install -y \
    postgresql-client \
    mysql-client \
    redis-tools \
    sqlite3

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create runner user
RUN useradd -m -s /bin/bash runner && \
    mkdir -p /home/runner/work && \
    chown -R runner:runner /home/runner

# Set environment variables
ENV CI=true
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

USER runner
WORKDIR /home/runner

CMD ["/bin/bash"]
``````

### Building and Using the Image

``````bash
# Build the image
docker build -t ubuntu-actions-runner:24.04 .

# Run the image
docker run -it --rm ubuntu-actions-runner:24.04

# Run with Docker socket (for Docker-in-Docker)
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock ubuntu-actions-runner:24.04

# Run with volume mount
docker run -it --rm -v $(pwd):/workspace -w /workspace ubuntu-actions-runner:24.04
``````

## Key Differences from Runner

Note aspects that cannot be perfectly replicated in a custom Docker image:

1. **GitHub Actions Context**: The runner includes GitHub Actions-specific environment variables (`GITHUB_WORKSPACE`, `GITHUB_SHA`, `GITHUB_REF`, etc.) and context that won't be available in a custom Docker image unless explicitly set.

2. **Pre-cached Dependencies**: The runner image has pre-cached dependencies (npm packages, Maven artifacts, etc.) for faster builds. Custom images start fresh each time unless you implement your own caching strategy.

3. **Service Configuration**: Some services (PostgreSQL, MySQL, MongoDB) are installed but not running by default in the runner. In a custom image, you need to handle service management separately.

4. **File System Layout**: The runner uses specific directory structures:
   - `/home/runner/work/{repo}/{repo}` for the workspace
   - `/home/runner/work/_temp` for temporary files
   - `/opt/hostedtoolcache` for language runtimes

5. **Hardware Resources**: GitHub-hosted runners have specific resource allocations:
   - 2-core CPU (7 GB RAM) for Linux runners
   - ~14 GB of SSD disk space for GITHUB_WORKSPACE
   - ~70-90 GB total disk space

6. **Network Configuration**: GitHub-hosted runners have unrestricted internet access and can reach GitHub services without additional configuration.

7. **Tool Versions**: Tool versions are updated regularly by GitHub. The runner image is rebuilt weekly with the latest stable versions of most tools.

8. **Browser Support**: The runner includes Chrome, Firefox, and Edge with WebDriver support. Replicating this in Docker requires additional configuration for headless operation.

9. **Authentication**: The runner has built-in authentication via `GITHUB_TOKEN`. Custom images need separate authentication setup.

10. **Immutability**: Each GitHub-hosted runner job starts with a fresh image. Custom Docker images need to ensure proper cleanup between runs.

## Optimization Tips

### For Faster Builds

1. **Use Multi-stage Builds**: Separate build dependencies from runtime dependencies
2. **Layer Caching**: Order Dockerfile instructions from least to most frequently changing
3. **Minimize Layers**: Combine related RUN commands with `&&`
4. **Clean Up**: Remove package manager caches and temporary files

``````dockerfile
# Example of optimized layer
RUN apt-get update && apt-get install -y \
    package1 \
    package2 \
    package3 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
``````

### For Smaller Images

1. **Use .dockerignore**: Exclude unnecessary files from build context
2. **Remove Build Tools**: Uninstall build dependencies after compilation
3. **Use Alpine-based Images**: For specific tools (not for the base Ubuntu image)

### For Security

1. **Run as Non-root**: Use the `runner` user for operations
2. **Update Regularly**: Rebuild images frequently to get security updates
3. **Scan for Vulnerabilities**: Use `docker scan` or Trivy
4. **Minimize Attack Surface**: Only install necessary tools

## Maintenance Notes

- The runner image is updated weekly by GitHub (usually on Sundays)
- Major version updates (e.g., Ubuntu 22.04 → 24.04) are announced months in advance
- Check the [actions/runner-images](https://github.com/actions/runner-images) repository for:
  - Release notes
  - Breaking changes
  - Deprecation notices
  - Version updates
- The `ubuntu-latest` label is updated to point to new Ubuntu versions approximately 6 months after the new Ubuntu LTS release
- Monitor the [GitHub Actions changelog](https://github.blog/changelog/label/actions/) for runner updates

## Troubleshooting

### Common Issues

**Issue**: Docker-in-Docker not working
- **Solution**: Mount the Docker socket: `-v /var/run/docker.sock:/var/run/docker.sock`

**Issue**: Permission denied errors
- **Solution**: Ensure the `runner` user has appropriate permissions or run as root temporarily

**Issue**: Tool not found
- **Solution**: Verify the tool is installed and PATH is correctly set

**Issue**: Out of disk space
- **Solution**: GitHub runners have ~14 GB for workspace. Clean up intermediate files

**Issue**: Network timeouts
- **Solution**: Check firewall rules and DNS configuration

### Debugging Tips

``````bash
# Check installed tool versions
node --version
python3 --version
docker --version

# Verify PATH
echo $PATH

# Check disk space
df -h

# Check memory
free -h

# List installed packages
dpkg -l | grep package-name

# Check service status (if systemd available)
systemctl status docker
``````

## References

- **Runner Images Repository**: https://github.com/actions/runner-images
- **Ubuntu 24.04 Documentation**: https://github.com/actions/runner-images/blob/main/images/ubuntu/Ubuntu2404-Readme.md
- **GitHub Actions Documentation**: https://docs.github.com/en/actions
- **Docker Documentation**: https://docs.docker.com/
- **Ubuntu Documentation**: https://ubuntu.com/server/docs
- **GitHub Actions Changelog**: https://github.blog/changelog/label/actions/
- **Runner Images Releases**: https://github.com/actions/runner-images/releases

## Version History

| Date | Ubuntu Version | Notable Changes |
|------|----------------|-----------------|
| 2026-02-07 | 24.04 LTS | Initial documentation created |
| TBD | 24.04 LTS | Future updates will be tracked here |

---

*This document is maintained as part of the gh-aw repository to help developers understand and replicate the GitHub Actions runner environment.*
