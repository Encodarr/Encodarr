<p align="center">
  <a href="https://github.com/Transfigurr/Transfigurr">
    <img src="https://raw.githubusercontent.com/Transfigurr/Transfigurr/main/frontend/public/transfigurr.svg" alt="Logo" width="150" height="150">
  </a>

  <h3 align="center"><a href="https://github.com/transfigurr/transfigurr">Transfigurr</a></h3>

  <p align="center">
    A modern codec management and automation tool
    <br/>
    <br/>
    <a href="https://www.transfigurr.media">View Demo</a>
    .
    <a href="https://github.com/Transfigurr/Transfigurr/issues">Report Bug</a>
    .
    <a href="https://github.com/Transfigurr/Transfigurr/issues">Request Feature</a>
  </p>
</p>

[![GitHub Stars](https://img.shields.io/github/stars/transfigurr/transfigurr.svg?color=9C27B0&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github)](https://github.com/transfigurr/transfigurr)
[![GitHub Release](https://img.shields.io/github/release/transfigurr/transfigurr.svg?include_prereleases&color=9C27B0&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github)](https://github.com/transfigurr/transfigurr/releases)
[![GitHub Package Repository](https://img.shields.io/static/v1.svg?color=9C27B0&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=transfigurr&message=GitHub%20Package&logo=github)](https://github.com/transfigurr/transfigurr/packages)
[![Docker Pulls](https://img.shields.io/docker/pulls/transfigurr/transfigurr.svg?color=9C27B0&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=pulls&logo=docker)](https://hub.docker.com/r/transfigurr/transfigurr)
[![Docker Stars](https://img.shields.io/docker/stars/transfigurr/transfigurr.svg?color=9C27B0&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=stars&logo=docker)](https://hub.docker.com/r/transfigurr/transfigurr)
[![GitHub Build Workflow Status](https://img.shields.io/github/actions/workflow/status/transfigurr/transfigurr/docker.yml?branch=main&style=for-the-badge&jobUrl=https%3A%2F%2Fci.transfigurr.com%2Fjob%2FDocker-Pipeline-Builders%2Fjob%2Ftransfigurr%2Fjob%2Fmaster%2F&logo=github)](https://ci.transfigurr.com/job/Docker-Pipeline-Builders/job/transfigurr/job/master/)

[Transfigurr](https://transfigurr.media/) is a modern media management tool specifically designed for automating codec encodings using FFmpeg. It streamlines the process of video file conversion according to custom ffmpeg profiles, making it an essential tool for media enthusiasts and professionals alike.

## Supported Architectures

Docker manifest is utilized for multi-platform awareness. More information is available from docker [here](https://distribution.github.io/distribution/spec/manifest-v2-2/#manifest-list).

Simply pulling `transfigurr/transfigurr:latest` should retrieve the correct image for your arch.

The architectures supported by this image are:

| Architecture | Available | Platform Flag |
| :----------: | :-------: | ------------- |
|    x86-64    |    ✅     | linux/amd64   |
|    arm64     |    ✅     | linux/arm64   |
|    armhf     |    ❌     |               |

<strong>Advanced:</strong> You can also pull an arch manually via the platform flag.

```
docker pull --platform linux/amd64 transfigurr/transfigurr:latest
docker pull --platform linux/arm64 transfigurr/transfigurr:latest
```

## Version Tags

This image provides various versions that are available via tags. Please read the descriptions carefully and exercise caution when using unstable or development tags.

|   Tag   | Available | Description          |
| :-----: | :-------: | -------------------- |
| latest  |    ✅     | Stable releases      |
| develop |    ✅     | Development releases |

## Application Setup

Access the webui at `<your-ip>:7889`, for more information check out [Transfigurr](https://transfigurr.media/).

### Media folders

The `/series`, `/movies` and `/transcode` have been set as **_optional paths_**, this is because it is the easiest way to get started. While easy to use, it has some drawbacks. Mainly losing the ability to hardlink (TL;DR a way for a file to exist in multiple places on the same file system while only consuming one file worth of space), or atomic move (TL;DR instant file moves, rather than copy+delete) files while processing content.

Use the optional paths if you dont understand, or dont want hardlinks/atomic moves.

The folks over at servarr.com wrote a good [write-up](https://wiki.servarr.com/Docker_Guide#Consistent_and_well_planned_paths) on how to get started with this.

## Usage

To help you get started creating a container from this image you can either use docker-compose or the docker cli.

### docker-compose (recommended)

```yaml
---
services:
  transfigurr:
    image: transfigurr/transfigurr:latest
    container_name: transfigurr
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - /path/to/data:/config
      - /path/to/series:/series #optional
      - /path/to/movies:/movies #optional
      - /path/to/transcode/transcode #optional
    ports:
      - 7889:7889
    restart: unless-stopped
```

### docker cli

```bash
docker run -d \
  --name=transfigurr \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Etc/UTC \
  -p 7889:7889 \
  -v /path/to/data:/config \
  -v /path/to/series:/series #optional \
  -v /path/to/movies:/movies #optional \
  -v /path/to/transcode:/transcode #optional \
  --restart unless-stopped \
  transfigurr/transfigurr:latest
```

## Parameters

Containers are configured using parameters passed at runtime (such as those above). These parameters are separated by a colon and indicate `<external>:<internal>` respectively. For example, `-p 8080:80` would expose port `80` from inside the container to be accessible from the host's IP on port `8080` outside the container.

|    Parameter    | Function                                                                                                       |
| :-------------: | -------------------------------------------------------------------------------------------------------------- |
|    `-p 7889`    | The port for the Transfigurr webinterface                                                                      |
| `-e PUID=1000`  | for UserID - see below for explanation                                                                         |
| `-e PGID=1000`  | for GroupID - see below for explanation                                                                        |
| `-e TZ=Etc/UTC` | specify a timezone to use, see this [list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List). |
|  `-v /config`   | Database and transfigurr configs                                                                               |
|  `-v /series`   | Location of series library on disk (See note in Application setup)                                             |
|  `-v /movies`   | Location of movies library on disk (See note in Application setup)                                             |
| `-v /transcode` | Location of transcode library on disk (See note in Application setup)                                          |

## Environment variables from files (Docker secrets)

You can set any environment variable from a file by using a special prepend `FILE__`.

As an example:

```bash
-e FILE__MYVAR=/run/secrets/mysecretvariable
```

Will set the environment variable `MYVAR` based on the contents of the `/run/secrets/mysecretvariable` file.

## User / Group Identifiers

When using volumes (`-v` flags), permissions issues can arise between the host OS and the container, we avoid this issue by allowing you to specify the user `PUID` and group `PGID`.

Ensure any volume directories on the host are owned by the same user you specify and any permissions issues will vanish like magic.

In this instance `PUID=1000` and `PGID=1000`, to find yours use `id your_user` as below:

```bash
id your_user
```

Example output:

```text
uid=1000(your_user) gid=1000(your_user) groups=1000(your_user)
```

## Support Info

- Shell access whilst the container is running:

  ```bash
  docker exec -it transfigurr /bin/bash
  ```

- To monitor the logs of the container in realtime:

  ```bash
  docker logs -f transfigurr
  ```

- Container version number:

  ```bash
  docker inspect -f '{{ index .Config.Labels "build_version" }}' transfigurr
  ```

- Image version number:

  ```bash
  docker inspect -f '{{ index .Config.Labels "build_version" }}' transfigurr/transfigurr:latest
  ```

## Updating Info

Below are the instructions for updating containers:

### Via Docker Compose

- Update images:

  - All images:

    ```bash
    docker-compose pull
    ```

  - Single image:

    ```bash
    docker-compose pull transfigurr
    ```

- Update containers:

  - All containers:

    ```bash
    docker-compose up -d
    ```

  - Single container:

    ```bash
    docker-compose up -d transfigurr
    ```

- You can also remove the old dangling images:

  ```bash
  docker image prune
  ```

### Via Docker Run

- Update the image:

  ```bash
  docker pull transfigurr/transfigurr:latest
  ```

- Stop the running container:

  ```bash
  docker stop transfigurr
  ```

- Delete the container:

  ```bash
  docker rm transfigurr
  ```

- Recreate a new container with the same docker run parameters as instructed above (if mapped correctly to a host folder, your `/config` folder and settings will be preserved)
- You can also remove the old dangling images:

  ```bash
  docker image prune
  ```

## Building locally

If you want to make local modifications to these images for development purposes or just to customize the logic:

```bash
git clone https://github.com/transfigurr/transfigurr.git
cd transfigurr
docker build \
  --no-cache \
  --pull \
  -t transfigurr/transfigurr:latest .
```
