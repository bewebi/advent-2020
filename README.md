# advent-2020
My [Advent of Code 2020](https://adventofcode.com/) solutions!

## Usage

Each day's solution can be found in a subdirectory.

The go executable can be built from the respective subdirectory by running

```bash
$ go build
```

and then run like

```bash
$ ./Day01 input.txt
2020/12/03 21:07:57 Inputs [933 1087] sum to 2020! Their product is 1014171
```

This will log the answer for Part 1 in a user-friendly message. Each executable accepts a flag that will result in the answer for Part 2 of the day's problem.

For example for Day 1:

```bash
$ ./Day01 -three-sum input.txt
2020/12/03 21:08:47 Inputs [1395 59 566] sum to 2020! Their product is 46584630
```

Run the executable with the `-help` flag to get a readout of the available flags:

```bash
$ ./Day01 -help
Usage of ./Day01:
  -target int
    	target sum (default 2020)
  -three-sum
    	find three input lines that sum to target rather than two
```

If improperly formatted input is given the program will log a fatal error.