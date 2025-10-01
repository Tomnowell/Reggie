# Reggie

A (mostly) recursive regex matcher in Go

## Use:

Example Use:

```bash
echo -n 'I see 1 cat' | ./reggie.sh -E '^I see \d+ (cat|dog)s?$'
```
### Why?

I know the World doesn't need another regex matcher. But this was a fun experiment and a learning tool for me to understand string and token parsing, acutally use
recursion for a "real world" problem, and get to know some basic Go syntax.

I'm definitely not convinced I chose the best way to tackle this problem. If you have any comments or advice, please do let me know.

## Credits:
While this code was written by me. I probably would never have started, let alone completed(sic) it without help from the following sources:

This regex matcher is heavily based on the example C code by Rob Pike found [here](https://www.cs.princeton.edu/courses/archive/spr09/cos333/beautiful.html)
I took his cryptic, yet beautiful code and made it ugly and a little more functional.

I created most of this matcher testing against the 'Build your own Grep' challenge on [CodeCrafters.](https://codecrafters.io) 
I enjoy the lack of handholding and level of freedom they provide with their challenges!
