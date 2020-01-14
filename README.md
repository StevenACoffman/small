# Get Small with Docker
### Create the smallest secured golang docker image

Restricting what's in your runtime container to precisely what's necessary for your app is a best practice employed by Google and other tech giants that have used containers in production for many years. It improves the signal to noise of scanners (e.g. CVE) and reduces the burden of establishing provenance to just what you need.

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
test                           - Runs go test with default values```
```

### Quickstart 

```
make build && make run
```

### Docker Image Size Comparison

| **Builder Stage** | **Final Stage** | | **Final Image Size** |
|---|---|---|
| `golang:1.13.6-alpine3.11` | `scratch` | **6.84MB**  |
| `golang:1.13.6-buster` | `gcr.io/distroless/static:nonroot` | **7.27MB** |

Alpine uses the musl library, and distroless uses glibc library. 
If you are using libraries that require cgo, sometimes they don't work well with musl.

**Alpine** is basically busybox linux with a package manager.
**Distroless** is basically debian _without_ a package manager.

Distroless has [different tags for base images](https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base?gcrImageListsize=10) not always mentioned in documentation :

+ latest
+ debug
+ nonroot
+ debug-nonroot

### Credit Where It Is Due

Alpine Docker container inspired by [this excellent article](https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324)