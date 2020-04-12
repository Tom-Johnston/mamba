# Mamba

Mamba is a set of Go packages for studying combinatorics. The main package is `graph` which provides data structures and functions for investigating small graphs. The API is not finalised and may change.

## Packages

- `comb`: Compute binomial coefficients and rank/unrank combinations.
- `dawg`: Create and search a directed acyclic word graph.
- `disjoint`: Create and manipulate a disjoint set data structure.
- `graph`: Create, manipulate and compute properties of small graphs.
  - `graph/search`: Generate all non-isomorphic graphs on n vertices (for very small values of n). It may be useful to copy and modify this code to search for graphs with certain properties.
- `ints`: Helper functions on `[]int`. Mostly a small subset of functions from the standard library's `byte` package translated to work on`[]int` instead.
- `itertools`: Iterate over permutations, combinations and set partitions.
-  `sortints`: Implements a set of `int` elements by storing them in a slice in ascending order.
- `tsp`: Output a travelling salesman problem in TSPLIB format for use in external solvers.
