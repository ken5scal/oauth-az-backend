# oauth-az-backend

# Package Structure

This repository is preferring `Layered Architecture` over `Flat Package Archtecture`, which many Go projects employ, for the following reasons.

## Easy to continue 
I build this for free time meaning there are times I cannot commit on this.
Layered architecture makes me easy to comeback because it defines specific responsibilities for each layer.
Maybe I can handle flat package architecture, but it requires me writing a clear and readable documents for me in the future.

## Pre-requisites for Flat Package architecture
Flat Package seems to require separating responsibility in each packages.
But, that's what exactly layered architecture achieves. 
Meaning flat package itself requires me well understanding the layered package.

## Not OSS
I'm committing this project not for OSS project; but more like for educational purpose.
If this were an OSS project, flat package architecture is preferred for other contributors.

# The Default Parameters

## Server environment
* running env: docker container
* port: 8080
* config file location: /etc/oauth-az/config.toml

One can customize parameters by adding `--build-args` in docker build command 

## App Environment
* debug: false 

* Build Using Docker
```
% docker build -f Dockerfile -t oauth-az-back-dev .
% docker run -it --rm -p 8080:8080 --name oauth-az-back oauth-az-back-dev:latest
```

* Build Locally
```
% go mod init
% vgo build
% vgo mod tidy
```

## DB Envinronment
```
% mkdir db
% sudo chmod -R 777 db
% docker run --name mysql --restart always -v $(pwd)/db:/var/lib/mysql -v $(pwd)/db/config:/etc/mysql/conf.d -e MYSQL_ROOT_PASSWORD=password -p 3306:3306 -d mysql
% docker exec -it mysql /bin/bash
root@0caa697a9dcc:/# mysql -u root -p
mysql >  CREATE USER 'ken5scal' IDENTIFIED WITH mysql_native_password BY 'password';
mysql > select User, Plugin from mysql.user;
```

# ToDo
* [ ] Implement OAuth
* [ ] Makefile
* [ ] Separate debug/prod envs