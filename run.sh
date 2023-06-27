kind delete cluster
kind create cluster
make manifests
make install
make run & kubectl apply -f ./config/samples
kubectl get spireserver