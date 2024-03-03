% Docker Debug Beta
% Mohammad-Ali A'râbi
% 2024-01-29

Docker Debug is a new feature available from Docker Desktop 4.27.0 that allows you to attach a shell to Docker container or image, even though they have no shell available.
The feature helps engineers to debug their containers and images without having to install additional tools in them. Docker Debug comes with a pre-installed set of tools, like editors, ping, etc.,
but you can also install your own tools.

# How to Use Docker Debug

To use Docker Debug, you need to have Docker Desktop 4.27.0 or later installed. Also, the beta version is only available for Pro subscribers.

To make sure the feature is available, run the following command:

```bash
docker debug --help
```

## Create an Image with No Shell

Let's first create a Docker image with no shell available. What follows is a Dockerfile based on `scratch` image. In the first stage, we create a binary file called `hello` and copy it to the second stage.

Create a Go file called `main.go` with the following content:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello, Wowlrd!")
		time.Sleep(5 * time.Second)
	}
}
```

It will print `Hello, Wowlrd!` every 5 seconds.
Now, create a Dockerfile with the following content:

```Dockerfile
FROM golang:1.17.2-alpine3.14 AS builder

WORKDIR /app

COPY main.go .

RUN go build -o hello main.go

# ------------------------------
FROM scratch

COPY --from=builder /app/hello /

CMD ["/hello"]
```

Build the image:

```bash
docker build -t go-hello .
```

Let's run the image:

```bash
docker run --rm --name hello go-hello
```

Well, I tagged the image also as `aerabi/go-hello` and pushed it to Docker Hub, in the case you don't want to build it yourself:
    
```bash
docker run --rm --name hello aerabi/go-hello
```

# Attach a Shell to a Container

So, we have an image with no shell available. It's not running in a container named `hello`.
Let's try to attach a shell to it using `docker exec`:

```bash
docker exec -it hello sh
```

It will fail with the following error:

```
OCI runtime exec failed: exec failed: unable to start container process: exec: "sh": executable file not found in $PATH: unknown
```

In short, there is no shell available in the image. Let's try to attach a shell using Docker Debug:

```bash
docker debug hello
```

The output would be something like this:

```
         ▄                                                                                                                                                                  
     ▄ ▄ ▄  ▀▄▀                                                                                                                                                             
   ▄ ▄ ▄ ▄ ▄▇▀  █▀▄ █▀█ █▀▀ █▄▀ █▀▀ █▀█                                                                                                                                     
  ▀████████▀    █▄▀ █▄█ █▄▄ █ █ ██▄ █▀▄                                                                                                                                     
   ▀█████▀                        DEBUG                                                                                                                                     
                                                                                                                                                                            
Builtin commands:                                                                                                                                                           
- install [tool1] [tool2] ...    Add Nix packages from: https://search.nixos.org/packages                                                                                   
- uninstall [tool1] [tool2] ...  Uninstall NixOS package(s).                                                                                                                
- entrypoint                     Print/lint/run the entrypoint.                                                                                                             
- builtins                       Show builtin commands.                                                                                                                     
                                                                                                                                                                            
Checks:                                                                                                                                                                     
✓ distro:            unknown/scratch                                                                                                                                        
✓ entrypoint linter: no errors (run 'entrypoint' for details)                                                                                                               
                                                                                                                                                                            
This is an attach shell, i.e.:                                                                                                                                              
- Any changes to the container filesystem are visible to the container directly.                                                                                            
- The /nix directory is invisible to the actual container.                                                                                                                  
                                                                                                                                                      Version: 0.0.22 (BETA)
root@7ce22feb81db / [hello]
docker > 
```

There are some things to unpack:

- The first line is a cool ASCII art. It's a whale with a shell.
- There are some builtin commands available. We will talk about them later.
- The description says "Nix" a lot. It's because Docker Debug is based on NixOS. It's a Linux distribution based on Nix package manager.
- The last line is the prompt. It's a shell prompt. You can run any command you want.

Let's try to run `ls`:

```bash
docker > ls
```

And here is the output:

```
dev  etc  hello  nix  proc  sys
```

You can see that there is a `hello` file. It's the binary file we created in the first stage of the Dockerfile.
And there are some other directories, like `nix`.

Let's try doing `netstat`:

```bash
docker > netstat
```

It will say it's not available. So let's install it:

```bash
docker > install net-tools
docker > netstat
```

As expected, there are no connections. On the container there are other tools already available like ping and curl.

# Attach a Shell to a Docker Image

Let's use another shell and stop the container:

```bash
docker stop hello
```

Now, there is no `hello` container running, but we have an image named `go-hello`.
Let's attach a shell to it:

```bash
docker debug go-hello
```

There is a warning saying:

```
Note: This is a sandbox shell. All changes will not affect the actual image.
```

This means a container is created from the image, and the shell is attached to it.
It's specially useful when you want to debug an image that is not running in a container.

Let's exit and run the container again:

```bash
docker run --rm -d --name hello go-hello
docker debug go-hello
```

Let's write something into the container filesystem:

```bash
docker > echo "Hello, World!" > hello.txt
```

You could also use Vim to create the file, as it's available in the environment.

Now, let's exit Docker Debug and commit the changes:

```bash
docker > exit
docker commit hello go-hello:2
```

Time to look into the newly created image:

```bash
docker debug go-hello:2
```

And if you look into the filesystem using an `ls` command, drum roll please:

```
dev  etc  hello  hello.txt  nix  proc  sys
```

It's there! The file we created in the container is now available in the image.
This means your changes to the container are actually persisted. How about the Nix packages we installed?

```bash
docker > netstat
```

It's there! But not because it was persisted in the container. The container didn't have a Nix ecosystem to begin with.
It's because once you install a Nix package on your Docker Debug shell, it will become available in all other Docker Debug shells.
This is actually desirable, because you don't want to install the same package over and over again.

## Attach a Shell to a Stopped Container

Let's run a container that stops immediately:

```bash
docker run --name hello hello-world
```

This is a special image that just prints "Hello from Docker!" and stops.
It's obviously impossible to attach a shell to it using `docker exec`, because it's not running.
But you can attach a shell to it using Docker Debug:

```bash
docker debug hello
```

This is specifically useful when you want to debug a container that stops immediately, e.g. if it shouldn't.
So, it's also different from looking into the image, because you can see the container filesystem and probably some logs.

# How Docker Debug Works

There exists a concept of "parked containers". It basically means that to get access to a container A, you create a new container B and mount the filesystem of A into B, along with some other things like network, process, etc.
In Linux, everything is a file. Your running processes, network interfaces, etc. are all files in the `/proc` directory. So, you can mount the `/proc` directory of A into B and see the processes of A.
Docker Debug is based on this concept. It's parks a NexOS container and gives you a shell to it.

Nix package manager is used for a few reasons:

- It's one of the most complete package managers.
- It allows to have different versions of the same package installed at the same time. It uses hashes of the package to differentiate them, so there is no conflict.

# Final Thoughts

Docker Debug is still in beta, so it might change in the future. The current version requires you to have a Pro subscription of Docker Desktop.
It's a very useful tool for debugging your containers and images. It's especially useful when you have a container that stops immediately, or an image with no shell available.

A Kubernetes version of Docker Debug is also in the works. It will allow you to attach a shell to a pod.

This article is based on the available Docker Debug CLI bundled into Docker Desktop, the official documentation, and a talk given by Johannes Großmann, the author of Docker Debug, at Docker Meetup in Freiburg, Germany.

I hope you find this article useful. If you have any questions or comments, please let me know. I'm always happy to help. 
I also write about Docker and git regularly. If you are interested,

- you can follow me on Twitter: [@MohammadAliEN](https://twitter.com/MohammadAliEN)
- or subscribe to my Medium blog: [https://medium.com/@aerabi](https://medium.com/@aerabi)
