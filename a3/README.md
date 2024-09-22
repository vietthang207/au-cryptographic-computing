# Asignment 3

go run a3.go

## Some explanation

I used a datastructure to represent the circuit which consist of 3 arrays (gates stores gate types, firstFanins and secondFanins store the 1 or 2 input to that gate). The fanin arrays store either an index or a constant. The exact convention is commented on the code.

Self-reflection: my usage of currentWire variable to keep track of the current state is somewhat convoluted and hard to follow (which means it is prone to error). It should be much easier if I have used another variable to keep track of the state (i.e. Sending, Receiving, ...)
