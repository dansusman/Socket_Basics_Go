# Socket Basics

This repo holds the source code for Socket Basics, the first project in my Networks and Distributed Systems course at NEU. It is written in Golang.

## High-level Approach

Given the protocol description, I aimed to break the problem into subtasks
initially. Knowing the server handles two of the four types of messages, the
main challenge was creating a client that supports writing the other two types
of messages (HELLO's and COUNT's) and reading the server-handled FIND's and
BYE's.

On top of this, there was the subtask of reading in command line arguments and
parsing for optional flags -p and -s. This was where I started, since it sounded
the simplest and most familiar in terms of overall concepts/implementations.

With just a few Google searches, I discovered the flag library for Golang, which
allows a developer to easily work with command line arguments, optional flags,
etc., so I made use of it. Everything was straightforward after finding the
library. Alternatively, I could have used argparse library (inspired heavily by
the Python argparse library) in Go, but I figured the basic flags library was
plenty good.

After the command line arguments are parsed and handled correctly, the client
sends a HELLO message via net.Dial or tls.Dial depending on the -s flag. It
talks on ports 27993 or 27994 based on the supplied or not supplied -p flag.

Using the net and crypto/tls packages seemed like a no-brainer, since they are
the go-to packages for these types of tasks in Golang. Luckily, I discovered
this Monday that *tls.Conn is actually just a specific net.Conn, meaning a
simple cast from *tls.Conn to net.Conn made abstracting the code easy.

I handle reading responses from the server with a bufio.NewReader, and, since we
know a well-formed response will always end with '\n', I utilize the
ReadString('\n') function. Barring some weird bugs I introduced in my error
handling, reading was somewhat straightforward to implement as well.

Writing to the server was as simple as calling the Write() of net.Conn. Easy!

Apart from that, the main function takes the subtasks and paints the full
picture: read the inputs, write HELLO, read FIND, verify the response is valid,
count the characters, write COUNT, read FIND, verify the response is valid,
loop, ..., read BYE, print secret flag.

## Problems Faced

- I set out to learn Golang this summer, but didn't get much further than
  fmt.Println("Hello, world!"), to be honest. So, writing this project in Golang
  was certainly a bigger challenge than using Python or something else I'm
  familiar with.

- I ran into a bunch of minor bugs like writing `ex_string COUNT {number} \n`
  instead of `ex_string COUNT {number}\n` (a single space between the number and
  the newline caused an hour or two of a headache). Also, setting up the initial
  connection to the server was a bit confusing when I first started the project,
  since I had never worked with TCP or Go.

- Until I read the docs for ReadLine(), I didn't understand what isPrefix was
  used for, so I just ignored it for a day of programming. That made for some
  useless reading of FIND messages. I stuck with it and implemented a version of
  my helper method that used ReadLine(). Then I found ReadString('\n') which
  simplified things further still.

## Testing Code

I tested the code on my own machine for the most part. My first step was writing
a Hello.go program in my Khoury Linux machine to ensure that Golang is indeed
supported up there, but after that most of my testing was from my laptop to
proj1.3700.network. My process was `go build client.go`, then running the
./client executable with various command line arguments to check that everything
works as expected. Overall, I tested as I developed.

Since the functionality was pretty clear to me, I opted out of writing tests
before writing code. At some points in the process I regretted that a little.
