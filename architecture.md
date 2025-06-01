I'm building a simple tui-based version of monkeytype in Go. I want to
replicate their core functionality for completing a typing test under timed
conditions, providing visual feedback on the user's progress with some type of
syntax coloring for words that are typed correctly versus mistakes. At the end
of a test the program should ask if the user wants to try again or see a log of
their stats, implying that summary statistics for each run should be stored.
There should also be escape conditions that let the user abort a typing test if
they desire to using escape. Standard SIGINT type functions should be supported
also (Ctrl+D to kill and Ctrl+C to interrupt).

Everything about this should live locally in the user's current working
directory. Statistics can be stored in a simple csv file, there is no need to
build complexity into this function yet. Any services run through the app should
have configurations provided here. Services should connect to the app through a
main method or server we configure as part of our app if there are services to
connect. However, I do not see a reason for us to need any external service
providers for this problem.

Part of building an app like this is going to be finding a corpus of words with
frequencies to generate our own sequences of words with. I think you could
probably generate the words yourself? If that seems reasonable, provide a list
of the 1000 word stems you think are most common in a `words.txt` file that we
can just load at the start of the program; also provide frequencies for words so
we can use them to generate the sequences we desire. I think those data should
be provided outside the program and loaded at startup instead of hard coding a
table of words into the program.

The first step is planning: I know nothing about Go so you will need to decide
what tools from the Go standard library to use and what dependencies to add. I
can add dependencies if necessary, and you should let me review dependencies
before adding them without my consent. I also don't know the first thing about
this kind of application, so you will need to propose all of the data structures
and algorithms that we are going to build on top of. Take a few seconds to
introspect about this before building.
