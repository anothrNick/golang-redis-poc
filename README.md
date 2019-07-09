A simple implementation of HTTP request rate limiting to a golang server with redis.

## Run Locally

```sh
# install dependencies
$ make install

# compile binary / build images
$ make build

# bring up containers
$ make up

# make request to server
$ curl -v localhost:5001
{"requests":30}

# ... after 30 requests within 1 minute
$ curl localhost:5001
Rate limit reached.

# shutdown containers
$ make down
```

## Reference

* https://github.com/go-redis/redis
* https://redislabs.com/redis-best-practices/basic-rate-limiting/
