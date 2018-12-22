# service that traverses a given root path and creates POST requests to given url with file info as body

To build image:
```bash
docker build {image_name} -t .
```
To Run it:
```bash
docker run  {image_name} -path {path} -address {address}
```
Or for help:

```bash
docker run {image_name} --help
```

