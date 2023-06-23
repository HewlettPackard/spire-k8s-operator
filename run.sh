kind delete cluster
kind create cluster
sleep 20
make manifests
make install
make run