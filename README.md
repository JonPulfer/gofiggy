# gofiggy

My first controller and learned a lot along the way.


I have a script to create the kind cluster that you run from the root of this 
repos like so: -

_NB: assumes docker is installed (running too otherwise you get an error) and also kind needs to be installed_

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

From this point you can create and update config maps and see what happens.
