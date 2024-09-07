# Assignment 1
Name: Dinh Viet Thang

## Compatibility between 8 blood types:

Both of the following rules need to be satisfied at the same time:

- Rule 1 (Rh rule): Positive types can receive from both positive and negative types. Negative types can receive from negative types only.

- Rule 2 (ABO rule): O can only receive from O, A can receive from A and O, B can receive from B and O, AB can receive from A, B and O.

## Functions

### Truth table:

| x  | y  | Compatible |
|----|----|------------|
| AB | AB | 1          |
| AB | A  | 1          |
| AB | B  | 1          |
| AB | O  | 1          |
| A  | AB | 0          |
| A  | A  | 1          |
| A  | B  | 0          |
| A  | O  | 1          |
| B  | AB | 0          |
| B  | A  | 0          |
| B  | B  | 1          |
| B  | O  | 1          |
| O  | AB | 0          |
| O  | A  | 0          |
| O  | B  | 0          |
| O  | O  | 1          |

### 3-bit encoding:

x is the recipient and y is the donor

x is represented by 3 bits x[0], x[1], x[2]

y is represented by 3 bits y[0], y[1], y[2]

- The first bits x[0], y[0] are 0 for negative and 1 for positive types respectively. Rule one is satisfied if: ```(NOT x[0]) OR y[0] = 1```

- The second bits x[1], y[1] represent antibody A. These bits are 0 for AB and A, and 1 for B and O.

- Similarly, the third bits represent antibody B. THese bits are 0 for AB and B, and 1 for A and O.

- Rule 2 is satisfied if: ```(NOT x[1]) OR y[1] = 1``` and ```(NOT x[2]) OR y[2] = 1```

Thus, we have a nice bitwise formula ```(~x) | y == 0b111```

