Adapted from en wikipedia org wiki Brainfuck Adding two values

++       Cell c0 to 2
> +++++  Cell c1 to 5

[        Start your loops with your cell pointer on the loop counter  c1 in our case
< +      Add 1 to c0
> -      Subtract 1 from c1
]        End your loops with the cell pointer on the loop counter

At this point our program has added 5 to 2 leaving 7 in c0 and 0 in c1
< .        Print out c0