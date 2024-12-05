container_name="pulsar_id_stats"
image_name="pulsar_id_stats"
image_version="1.1.11"

run_debug: build_go
	./${container_name} --log_level=debug --source_allow_insecure_connection=true

run_trace: build_go
	./${container_name} --log_level=trace --source_allow_insecure_connection=true

run_info: build_go
	./${container_name} --log_level=info --source_allow_insecure_connection=true

curl:
	curl localhost:7700/metrics

build_go:
	go build -C src/main -o ../../${container_name}

pprof: build_go
	./${container_name} --pprof_on=true --log_level=debug

launch_pprof:
	go tool pprof -http=:8081 pprof/202*.pprof

run_container: build_cache
	docker run -d --rm --net host \
    	-p 7700:7700 --name ${container_name} ${image_name}:${image_version}

build:
	echo "Building ${image_name}:${image_version} --no-cache"
	docker build -t ${image_name}:${image_version} . --no-cache

build_cache:
	echo "Building ${image_name}:${image_version} --with-cache"
	docker build -t ${image_name}:${image_version} .

docker_hub: build
	./../push_image_dockerhub.sh ${image_name} ${image_version}

start_pulsar:
	docker run -d --rm --name pulsar \
		-p 6650:6650 -p 8080:8080 apachepulsar/pulsar:latest bin/pulsar standalone

clean:
	rm pulsar_id_stats || \
	rm pprof/2023*