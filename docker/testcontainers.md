# Testcontainers with Docker Compose and Node.js

This guide demonstrates how to use Testcontainers for a Node.js environment. 

Testcontainers is a tool for creating containers from within your application. A good example to demonstrate its capabilities is to use it with a Node.js application.
An Express web server that needs a Redis instance would usually require the user to install Redis on their machine, or at least spin it up using Docker. With Testcontainers, you can do this from within your code.

Let's see how it works.

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

After cloning the repository, navigate to the root and then to the `testcontainers` directory. Then install dependencies:

```bash
npm install
```

At this point, make sure you hava Testcontainers Desktop running. Let's use the "embedded runtime" for this example.
On my machine, it's called "Containers running locally".

Now let's start the app:

```bash
npm run dev
```

The app checks if the `REDIS_URL` environment variable is set. If so, it connects to the specified Redis instance. If `REDIS_URL` is not set, it uses testcontainers to create a Redis instance using the Docker Compose file `redis.yaml`.
To verify that a Redis instance is running, you can use the following command:

```bash
docker ps
```

We can see that two containers are running: the Redis instance and the Testcontainers instance.
Let's try to send a request to the app:

```bash
curl localhost:3000
```

It says "You are visitor number 1". If we send another request:

```bash
curl localhost:3000
```

It says "You are visitor number 2". The numbers are stored in the Redis instance.

So, now let's take a deeper look into how the whole thing works.

## How it works

The command that we ran (`npm run dev`) starts the app. Checking the `package.json` file, we can see that it runs the `dev` script:

```json
"scripts": {
  "dev": "nodemon --exec node --loader ts-node/esm index.ts"
}
```

The `dev` script uses `nodemon` to watch for changes in the files and restart the app when a change is detected. The `--exec node --loader ts-node/esm index.ts` part tells `nodemon` to run the `index.ts` file using the `ts-node` loader.
So what's inside the `index.ts` file? A normal Express app, but it imports the content of `setupRedis.ts`, which is where the magic happens:

```typescript
const environment = await new testcontainers.DockerComposeEnvironment(".", "redis.yaml")
  .withWaitStrategy("redis", testcontainers.Wait.forLogMessage("Ready to accept connections"))
  .withNoRecreate()
  .up()
```

This code creates runs a Docker Compose environment using the `redis.yaml` file, with the wait strategy set to wait for the log message "Ready to accept connections". This message is printed by the Redis instance when it's ready to accept connections. The `withNoRecreate()` method tells Testcontainers not to recreate the containers if they already exist. The `up()` method starts the containers.

## Conclusion

So far, a server application should've relied on third-party means to prepare the environment. With Testcontainers, from within your code you can run Docker containers needed for your application to work.
The technology is called Testcontainers, as it's not an optimal solution for production. However, the naming shouldn't restrict your creativity.
