kind create cluster
sleep 20
make manifests
make install
make run
sleep 30
kind delete cluster