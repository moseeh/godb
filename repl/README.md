# REPL Package

The `repl` package provides a Read-Eval-Print Loop (REPL) for interacting with the `godb` database from the command line.

## Import

To use the `repl` package, import it as follows:

```go
import "godb/repl"
```

## Usage

The main component of the `repl` package is the `REPL` struct. You can create a new `REPL` instance using the `NewREPL` function, which takes an `io.Reader` as input (e.g., `os.Stdin`). The `Start` method begins the REPL loop.

**Example:**

To start a new REPL session that reads from standard input:

```go
package main

import (
	"godb/repl"
	"os"
)

func main() {
	r := repl.NewREPL(os.Stdin)
	r.Start()
}
```

When the REPL is running, you can enter SQL commands at the prompt:

```
godb> CREATE TABLE users (id INT PRIMARY KEY, name STRING NOT NULL);
✓ Table 'users' created successfully
godb> INSERT INTO users (id, name) VALUES (1, 'moses');
✓ 1 row inserted
godb> SELECT * FROM users;
id | name
---+------
1  | moses

1 row(s) returned.
```

## Components

### REPL Struct

The `REPL` struct manages the state of the REPL session, including the database instance and the input reader. It contains the main loop that reads user input, sends it to the parser, executes the resulting command, and prints the output.

### Printer Functions

The `printer.go` file provides helper functions for formatting and printing output to the console.

-   `PrintRows`: Formats and prints a slice of `engine.Row` in a user-friendly table format.
-   `PrintSuccess`: Prints a success message to the console.
-   `PrintError`: Prints an error message to the console.