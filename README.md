# advent-2020
My [Advent of Code 2020](https://adventofcode.com/) solutions!

## Rankings

```
      --------Part 1---------   --------Part 2--------
Day       Time    Rank  Score       Time   Rank  Score
 22   00:09:34    1002      0   01:21:33   2358      0
 21   00:32:54    1242      0   00:43:19   1214      0
 20   00:38:17     794      0   13:03:43   3817      0
 19   00:26:14     432      0   02:18:52   1939      0
 18   00:24:38    1077      0   00:45:41   1292      0
 17   00:36:15    1490      0   00:48:23   1622      0
 16   00:25:48    2925      0   00:59:13   1788      0
 15   00:18:48    1869      0   00:19:17    796      0
 14   00:58:40    5039      0   01:10:46   2901      0
 13   01:41:50    8271      0   02:17:38   3606      0
 12   00:20:35    2870      0   00:34:23   2005      0
 11   00:39:37    3370      0   01:04:29   3147      0
 10   00:10:08    2349      0   00:39:38   2003      0
  9   00:28:48    6670      0   01:06:54   7954      0
  8   00:34:38    7507      0   01:05:32   6785      0
  7   01:09:31    6077      0   01:14:56   4202      0
  6   00:09:49    3604      0   00:20:31   3772      0
  5   00:31:40    5686      0   00:36:09   4700      0
  4   00:38:24    7038      0   01:06:22   4788      0
  3   18:25:25   59841      0   18:52:13  57843      0
  2       >24h   80997      0       >24h  78540      0
  1       >24h  101832      0       >24h  95449      0
```

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

This will log the answer(s) in a user-friendly message. In cases where the default behavior does not output answers for both parts, the default behavior will output an answer for Part 1 and the user may specify a flag in order to get the answer for Part 2.

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

For most days, if improperly formatted input is given the program will log a fatal error.