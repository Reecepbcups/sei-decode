# juno-decoder

This just reduces overhead from the normal app so it ONLY decodes transactions. Used to speed up <https://github.com/Reecepbcups/cosmos-indexer>

```bash
make install

# Binaries have this, but this reduces overhead
# You can find these txs with https://rpc.network.com/block?height=7793647
juno-decode tx decode [amino-base64]

# Example in action
juno-decode tx decode --output json CrUBCrIBCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QSiQEKK2p1bm8xamc5ZjNkbnJzZzNkN25wNTRyanlxY3J6MHphZXcydjNyejVlNGQSP2p1bm8xNjRjZDZ3cGVwa3hwNG5kd2t2d2U5bnFwOWtyYzM3bHduNzB0N3h5M2g3amFhNHlwZWtzc3p1MnczdxoZeyJkaXN0cmlidXRlX3Jld2FyZHMiOnt9fRJpClIKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECJoHdaIdtGcDLsWxbEngOKq6GU6S3ykOXO8HeXt5H6EMSBAoCCAEYvtgDEhMKDQoFdWp1bm8SBDE2MDIQteMJGkAM8dn7/9tazy1CJGClB/GLP5TXzJIxUTcKHnP93kQKBGKBDqkxEUUbzzVtoNYLw/K0CMSzsMKibVQjVK6D7xfR
```
