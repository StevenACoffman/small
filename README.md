# Get Small with Distroless Docker and Go
### Create the smallest secured golang docker image

Restricting what's in your runtime container to precisely what's necessary for your app is a best practice employed by Google and other tech giants that have used containers in production for many years. Security scanners finding vulnerabilities (e.g. CVE) in unrelated baggage just adds tedious maintenance work. Skip all that and just maintain only what you need.

Running your containers as a non-root user prevents malicious code from gaining permissions in the container host and means that not just anyone who has pulled your container from the Docker Hub can gain access to everything on your server, for example.

```
✓ usage: make [target]

build                          - Build the application
clean                          - Cleans the binary
docker-build-no-cache          - Build the smallest secure golang docker image based on distroless static with no cache
docker-build                   - Build the smallest secure golang docker image based on distroless static
docker-push                    - Pushes the docker image to registry
docker-run                     - Run the smallest and secured golang docker image based on distroless static
help                           - Show help message
lint                           - Lint the application code for problems and nits
ls                             - List size docker images
run                            - Runs go run main.go
test                           - Runs go test with default values
```

`make run` will start a webserver that will listen and respond http://127.0.0.1:8080/health and everything else gives 404 Not Found.
### Quickstart 

```
make build && make run
```

#### [Docker Security Best Practices](https://snyk.io/blog/10-docker-image-security-best-practices/):
1. Prefer minimal base images
2. Least privileged user
3. Sign and verify images to mitigate MITM attacks
4. Find, fix and monitor for open source vulnerabilities
5. Don’t leak sensitive information to Docker images
6. Use fixed tags for immutability
7. Use COPY instead of ADD
8. Use metadata labels
9. Use multi-stage build for small and secure images
10. Use a [linter](https://github.com/hadolint/hadolint)

### Docker Image Size Comparison

| **Builder Stage** | **Final Stage** | **Final Image Size** |
|---|---|---|
| `golang:1.13.6-alpine3.11` | `scratch` | **6.84MB**  |
| `golang:1.13.6-buster` | `gcr.io/distroless/static:nonroot` | **7.27MB** |

Alpine uses the musl library, and [Distroless](https://github.com/GoogleContainerTools/distroless/tree/master/base) uses glibc library. 
If you are using libraries that require cgo, sometimes they don't work well with musl.

+ **Alpine** is basically busybox linux with a package manager.
+ **Distroless** is basically debian _without_ a package manager.

`gcr.io/distroless/static` contains:

* ca-certificates
* A /etc/passwd entry for a root user and nobody (unprivileged)
* A /tmp directory
* tzdata

The `Dockerfile.alpine` here shows how securing alpine-based docker builds is more complicated than with distroless.

BTW, Distroless has [different tags for base images](https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base?gcrImageListsize=10) not always mentioned in documentation :

+ latest
+ debug
+ nonroot
+ debug-nonroot

### Credit Where It Is Due

Alpine Docker container inspired by [this excellent article](https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324)

### Wanna get **real** small?
If you add upx to your builder stage, you can shrink the binary even more:
`upx --brute app`
This is is very slow to build and minutely impacts startup time and resources.

However, you can often fit the result on a floppy disk, so that's cool.

### Graceful shutdown
This is a pretty good writeup about how to do [Graceful shutdown](https://rafallorenz.com/go/handle-signals-to-graceful-shutdown-http-server/)