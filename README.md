# buildbot-dashboard [![Join the chat at https://gitter.im/ghophp/buildbot-dashboard](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/ghophp/buildbot-dashboard?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Build Status](https://travis-ci.org/ghophp/buildbot-dashboard.svg?branch=master)](https://travis-ci.org/ghophp/buildbot-dashboard) [![Coverage Status](https://coveralls.io/repos/github/ghophp/buildbot-dashboard/badge.svg?branch=master)](https://coveralls.io/github/ghophp/buildbot-dashboard?branch=master)
If you have a CI/CD setup, it is really important to keep a central display and follow the builds and process going on the CI tool. This project aims to be the dashboard that will allow you to do present better results at a central display with `buildbot`.

## Features
- Non-reload monitoring
- Enhanced UI with better visualization of the builders
- Easy usage with single command
- `Filter` options allow you to just show what matters
- Save arrangement of the dashboard
- `Zoom compliant` layout

## Preview
![Apache Board](./preview/preview_apache.gif?raw=true "Apache Board")

As we use gridster and a structured style, if your buildbot has a lot of builders, you can take advantage of the `browser zoom` to reduce the size of the blocks and so, use more space to move your grid.

![Apache Board Small](./preview/preview_apache_small.gif?raw=true "Apache Board Small")

## Running with Docker

```bash
$ docker run \
    -e "STATIC_PATH=/go/src/github.com/ghophp/buildbot-dashboard/static/" \
    -e "BUILDBOT_URL=https://ci.apache.org" \
    -e "PORT=8000" \
    -e "MARTINI_ENV=development" \
    -p 127.0.0.1:8000:8000/tcp ghophp/buildbot-dashboard
```

Or if you have clone/downloaded the repository you can run locally with:

```bash
$ make run.docker
```

This will build the image with the name `ghophp/buildbot-dashboard` and run it with the `.env.example` file if no `.env` file is present at the directory. To use other environment variables, please create a `.env` file next to the `.env.example` file.

## Running from Binary

First you have to build the project for you current architecture:

```bash
$ make build.local
$ make build.linux.armv8
$ make build.linux.armv7
$ make build.linux
$ make build.osx
$ make build.windows
```

This will yield a binary into the bin folder, which can be executed, `buildbot` is the only required flag, you must provide the base url of the running builbot.

```sh
$ ./buildbot_dashboard -h
-buildbot string
	buildbot url eg. http://10.0.0.1/
-filter string
	regex applied over the builder name
-invalidate int
	cache invalidate in seconds (default and min 5 minutes) (default 10)
-refresh int
	refresh rate in seconds (default and min 20 seconds) (default 20)
```

How to use with [runit](https://github.com/ghophp/buildbot-dashboard/wiki/runit).

## Running with go run

```sh
$ BUILDBOT_URL="https://ci.apache.org" PORT=8000 MARTINI_ENV=development make run.local
```

## Martini
As this project is build over `martini` please consider setting this `env variables` when deploy:
```sh
export PORT=3000
export MARTINI_ENV=production
```

## Clean Cache
By default the project will cache the builders to avoid the delay of reloading it every time. If you insert a new builder or remove one, you can force the reload of the cache in the UI, there is [a button on left side](https://github.com/ghophp/buildbot-dashboard/wiki/). This button will force the builders to be reloaded and the localStorage will also be cleaned.

---

Thanks to [nbari](https://github.com/nbari) and [paw](https://luckymarmot.com/paw) for helping and supporting open source.

![](./preview/paw.png)
