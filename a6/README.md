# Asignment 4

```
go run a6.go
```

I name the parameters as following:

- η: bit length of the secret key p

- N: number of integers in the public key (i.e. number of q_i and r_i)

- γ: bit length of each quotent q_i

- ρ: bit length of each 

The parameters need to satisfied that the noise are always less than p (i.e. output noise < min p).

The formulas for the output noise and min p is given in the code (it might be off by a few factor if my combinatorics math is incorrect, but I think the order of magnitude is right if we use big security params). 

I want to stay within 64 bit uint, so my params are quite small, most notably γ. In practice, I think we should have γ on par with η so that the underlying hardness assumption holds.

Regarding the choice of subset. I use a bitmask to randomize the subset. I think that any choice of subsets will satisfy correctness. However, if there are too few (less than N/4) or too many (more than 3N/4) elements are selected, then I resample the subset. I guess that the problem is too easy when the subset is sparse (for example, empty subset). I'm not sure if a nearly full subset makes the problem easy or not, so I remove them for the sake of symmetry. 

Run automated test. This fails with a small probability, so my math might be off a little bit.
```
bash test.sh
```
