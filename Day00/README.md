<h2 id="ex00">Exercise: Anscombe</h3>

So, let's say we have a bunch of integer numbers, strictly between -100000 and 100000. It may 
probably be a large set, so let's assume our application will read it from a standard input, 
separated by newlines. Right now let's think of four major statistical metrics that we can derive
from this data, so by default we can print all of them as a result, for example, like this:

```
Mean: 8.2
Median: 9.0
Mode: 3
SD: 4.35
```

The order and actual format doesn't really matter as long as we can understand which is which. 
A couple of notes, though:

1) Input data may or may not be sorted. You don't need to write your own sorting algorithm,
luckily, Go already has one in standard library and it works for integers.
2) Median is a middle number of a sorted sequence if its size is odd, and an average between
two middle ones if their count is even.
3) Mode is a number which is occurring most frequently, and if there are several, the smallest one
among those is returned. You may think of using some structure for storing number counts, and some
Go standard container may help.
4) You can use both population and regular standard deviation, whichever you prefer.
5) Calling someone "average" can be mean.

It will also make sense for user to be able to choose specifically, which of these four parameters
to print, so need to implement this as well. By default it's all of them, but there should be 
a way to specify whether print just one, two or three specific metrics out of four when running
the program (without recompilation). There is a built-in module in standard library that allows you 
to parse additional parameters.

<h2 id="rules-of-the-day" >Rules of the day</h2>

- You should only turn in `*.go` files and (in case of external dependencies) `go.mod` + `go.sum`
- Your program should accept a sequence of numbers separated by newlines to its standard input. One number is also a sequence.
- You may assume it should only work with integers (the output can be floats though, rounded to 2 decimal places).
- Nevertheless it should print a meaningful error message without runtime panic if fed some unexpected input, say, number out of bounds, letter characters or empty string.
- Your code for this task should be buildable with just `go build`