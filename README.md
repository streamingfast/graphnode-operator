# WHAT IS THIS

This is an operator to run some kubernetes jobs like running comparative tests on graph-nodes that require their own databases

# HOW THIS REPO WAS BOOTSTRAPPED

```
operator-sdk init --domain streamingfast.io --repo github.com/streamingfast/graphnode-operator
operator-sdk create api --group graphnode --version v1alpha1 --kind TestRun --resource --controller
```

# DEV FLOW

* run `make install run` with proper kubectl configuration
* deploy a 'testrun' object on your kubernetes (see config/samples)
* enjoy!

# HOW TO USE THIS

* Use the https://github.com/streamingfast/cloudsql-tarballer helper to generate a tarball of PGDATA from cloudSQL in a google storage bucket (ex: gs://example-bucket/tarballs)
* Run this operator on your k8s infrastructure (either with "make docker-build, make docker-push, make deploy IMG=..." or locally with proper kubectl permissions and `make install run`)
* copy and edit `config/samples/graphnode_v1alpha1_testrun.yaml` to match your environment
* deploy that testrun.yaml to run a job that will:
  1. start postgresql from the pgdata tarball
  2. compile graph-node from given github repo and branch
  3. run graph-node over that local database until a specific stop-block `GRAPH_STOP_BLOCK` is reached in one of the deployments with `GRAPH_DEBUG_POI_FILE` set to a local file
  4. uplaod the poi.csv file to google storage

# LIMITATIONS / REQUIREMENTS

* The tested graph-node must include `GRAPH_STOP_BLOCK` and `GRAPH_DEBUG_POI_FILE` handling
* The `GRAPH_STOP_BLOCK` mechanism is hacked and could be improved to require all working deployments to match, or include a failsafe timeout mechanism

