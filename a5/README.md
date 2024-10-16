# Asignment 4

```
go run a5.go
```

My prime testing algorithm is not very efficient, so I only used 26-bit prime (DDH_PARAM_K = 26)

I don't want to deal with overflow or big int, so, I only use 32-bit integer. Therefore, for the garbled circuit, I use PARAM_K = 15. For PRF, I just take the last 30 bits of the digest.

If garbled evaluation failed (most likely because there are multiple match when comparing the garbled table), my program will panic. With the current setting, it will happens with probability of about 10%.

My observation is that the collision chance increase if we decrease PARAM_K. If PARAM_K=12, it is about 50%.
