IMG_VERSION ?= 1.0.0

build:
	docker build -t go-sample-app:$(IMG_VERSION) .

run: build
	docker run --rm -e SOME_ENV=$$LALA -p 9090:9090 go-sample-app:$(IMG_VERSION)
