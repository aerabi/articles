# Testcontainers with Docker Compose and Node.js

This guide demonstrates how to use Testcontainers for a Node.js environment. 

Testcontainers is a tool for creating lightweight, throwaway instances of common databases or anything that can run in a Docker container. In addition to testing, you can leverage the Testcontainers [Docker Compose Feature](https://node.testcontainers.org/features/compose/) to perform **local development** in scenarios where you cannot install the Docker engine or need **cloud workloads** to offload resource-intensive tasks.

With Testcontainers Desktop, you can develop against [Testcontainers Cloud](https://testcontainers.com/cloud/) or the [Testcontainers embedded runtime](https://newsletter.testcontainers.com/announcements/adopt-testcontainers-desktop-as-your-container-runtime-early-access) to see Testcontainers capabilities.

## Prerequisites

- [Git](https://git-scm.com/downloads)
- [Node.js](https://nodejs.org/en/download/)
- [Docker](https://docs.docker.com/get-docker/)
- [Testcontainers Desktop](https://testcontainers.com/desktop/)

## Running the app

First clone the repository:

```bash
git clone git@github.com:docker/awesome-compose.git
```

After cloning the repository, navigate to the root and then to the `testcontainers` directory. You can directly start the project with Docker Compose:

```bash
docker compose up -d
```

It builds the Docker image and starts the container. In addition, it creates a Redis instance using the Docker Compose file `compose.yaml`.
