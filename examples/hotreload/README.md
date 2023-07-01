# Hot reload go application using PolyWatch
This is a development environment including an api service written in go which is run from source code
and will get recompiled whenever any `*.go` file under [api](api) directory get changed.

After go v1.20 there were some issues with source code mounting and ability to run go module maintenance commands & compile
while using basic docker volume: [fatal: 'origin' does not appear to be a git repository][iss-go-origin-git].
One solution was to mount using [lebokus/bindfs][git-bindfs] docker volume plugin. To install it:

```shell
docker plugin install lebokus/bindfs
```

When using [lebokus/bindfs][git-bindfs] docker volume plugin to mount source codes in docker container we'll have isolated files ownership:
* inside the **api** docker container all mounted files belong to **root** user
* outside the docker container (host) files belong to you
And be aware that [bindfs cannot sense fsnotify][iss-bindfs-7]. So we have to use polling method to watch for file changes

To bring up the environment:

```shell
docker-compose up
```

Visit api server on your browser:
* Local URL [http://api.hotreload.plw](http://api.hotreload.plw) if you're using [autodns][git-autodns]
* Find IP Address by running `docker container inspect --format='{{.NetworkSettings.Networks.hotreload_default.IPAddress}}' hotreload-api-1`

Then change source code of api server & check docker logs and refresh the page to notice that
the service has been recompiled (pid differs in response), and it includes the recent changes in source code.

[git-bindfs]: https://github.com/clecherbauer/docker-volume-bindfs
[git-autodns]: https://github.com/pouyanh/autodns
[iss-go-origin-git]: https://stackoverflow.com/questions/15637507/fatal-origin-does-not-appear-to-be-a-git-repository
[iss-bindfs-7]: https://github.com/mpartel/bindfs/issues/7
