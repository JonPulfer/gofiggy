# gofiggy

My first controller and learned a lot along the way.

I have a script to create the kind cluster the you run from the root of this 
repos like so: -

```bash
./bin/create_cluster.sh
kubectl cluster-info --context kind-kind
```

You can then build the gofiggy controller like: -

```bash
./bin/rebuild_gofiggy.sh
```

This will build the docker image and push it into the registry that is available
in the kind cluster.

To run the controller: -

```bash
kubectl create -f gofiggy.yaml
```

From this point you can create and update configmaps and see what happens.
