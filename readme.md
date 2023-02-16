    rpk container start -n 1
    rpk topic create product --brokers localhost:51826
    rpk container purge
    kcat -C -b localhost:50774 -t product -o beginning


    rpk container start -n 1
    rpk topic create product --brokers localhost:50774


    docker run  --network=host  -p 8080:8080 -e KAFKA_BROKERS=localhost:50774 docker.redpanda.com/vectorized/console:latest
    docker run  --network=host  -p 8081:8080 -e KAFKA_BROKERS=host.docker.internal:50774 docker.redpanda.com/vectorized/console:latest
