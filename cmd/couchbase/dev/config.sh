#!/bin/bash

retry() {
    for i in $(seq 1 5); do
        $1
        if [[ $? == 0 ]]; then
            return 0
        fi
        sleep 1
    done

    return 1
}

create-bucket(){
    for bucket in $BUCKETS; do
        echo "Creating $bucket bucket..."
        flush=0

        couchbase-cli bucket-create -c localhost -u $USERNAME -p $PASSWORD \
        --bucket=$bucket \
        --enable-flush=$flush \
        --bucket-type=couchbase \
        --bucket-ramsize=100 \
        --wait
    done
}

create-index(){
    for bucket in $BUCKETS; do
        cmd='CREATE PRIMARY INDEX ON `'$bucket'`'
        createOutput=$(cbq -u=$USERNAME -p=$PASSWORD --script="$cmd")
        

        if [[ $createOutput == *"\"status\": \"success\""* ]]; then
            echo "Index on $bucket successfully created"
        else
            echo "$bucket index creation FAILED"
            echo $createOutput
        fi
    done
}

add-apiclient-data() {
    cmd='insert into `api_client`
        VALUES ("admin",{ 
            "client_id":"admin", 
            "api_key":"admin", 
            "role":"admin"
            })
        RETURNING META().id as docid, *;'
    
    createOutput=$(cbq -u=$USERNAME -p=$PASSWORD --script="$cmd")
    
    if [[ $createOutput == *"\"status\": \"success\""* ]]; then
        echo "Data successfully inserted"
    else
        echo "Data insert FAILED"
        echo $createOutput
    fi
}

add-role-data() {
    cmd='insert into `service_roles` ( KEY, VALUE ) 
        VALUES
        (
            "admin",
            {"role_name":"admin", "allowed_paths":["/admin/*","/api/*"]}
        ), 
        VALUES (
            "api-client",
            {"role_name":"api-client", "allowed_paths":["/api/*"]}
        )
        RETURNING META().id as docid, *;'
    
    createOutput=$(cbq -u=$USERNAME -p=$PASSWORD --script="$cmd")
    
    if [[ $createOutput == *"\"status\": \"success\""* ]]; then
        echo "Data successfully inserted"
    else
        echo "Data insert FAILED"
        echo $createOutput
    fi
}

init-couchbase() {
    # wait for service to come up
    echo "Waiting for couchbase service..."
    
    until curl --output /dev/null --silent --head --fail http://localhost:8091; do
        echo "."
        sleep 1
    done
}

init-cluster(){
    # check if cluster is initialised
    initCheckOutput=$(curl --silent -u $USERNAME:$PASSWORD http://localhost:8091/pools/default)

    if [[ $initCheckOutput == *"unknown pool"* ]]; then
        # initialize cluster
        couchbase-cli cluster-init -c localhost \
            --cluster-username=$USERNAME \
            --cluster-password=$PASSWORD \
            --cluster-port=8091 \
            --services=data,index,query,fts \
            --cluster-ramsize=$MEMORY_QUOTA \
            --cluster-index-ramsize=$INDEX_MEMORY_QUOTA \
            --cluster-fts-ramsize=$FTS_MEMORY_QUOTA \
            --index-storage-setting=default

        echo "Couchbase cluster initialised"
    fi
}

main(){
    echo "Couchbase UI :8091"
    echo "Couchbase logs /opt/couchbase/var/lib/couchbase/logs"
    
    ./entrypoint.sh couchbase-server &

    init-couchbase
    init-cluster
    
    bucketCheckOutput=$(curl --silent -u $USERNAME:$PASSWORD http://localhost:8091/pools/default/buckets)
    
    if [[ $bucketCheckOutput == "[]" ]]; then
        retry create-bucket

        echo "Waiting for query service ..."
        sleep 5

        retry create-index
        retry add-apiclient-data
        retry add-role-data
    fi

    echo "Couchbase is ready :)"
    wait
}

main