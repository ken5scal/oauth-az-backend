# oauth-az-backend


# The Default Parameters

## Server environment
* running env: docker container
* port: 8080
* config file location: /etc/oauth-az/config.toml

One can customize parameters by adding `--build-args` in docker build command 

## App Environment
* debug: false 

# Build Using Docker
```
% docker build -f Dockerfile -t oauth-az-back-dev .
% docker run -it --rm -p 8080:8080 --name oauth-az-back oauth-az-back-dev:latest
```

# Build Locally
```
% go mod init
% vgo build
% vgo mod tidy
```

# ToDo
* [ ] Implement OAuth
* [ ] Makefile
* [ ] Separate debug/prod envs