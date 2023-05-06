cockroach start --advertise-addr 'localhost' --insecure  --join=roach1,roach2
echo Wait for servers to be up
sleep 5
cockroach init --insecure