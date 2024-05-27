# README

## What

I've created a simple CLI called _quoter_ using only the _flags_ package in Go, instead of something with batteries attached like _cobra_

## Usage

Just roughly

Main command: _quoter_ (or whatever name you compile the program to)

Sub-commands:

    - add: add a quote to an on-disk file

    - quote: retrieve a random quote from our on-disk file

Sub-command flags:
    - g: to specify the genre of the quote being added/retreived


## General structure

- Main command parses sub-commands and calls respective sub-command handler

- Sub-command has a _driver_ which handles side-effects of writing and calling functional components.

## Thoughts

- Functions for each sub-command are mostly structured in an easily testable, functional style. This talk about [Functional Core, Imperative Shell](https://news.ycombinator.com/item?id=18043058) is a good way to think about it. Good comments in the HN comment section as well.

- ^ (Personal opinion) No need to be too dogmatic about this though, I feel a good way to think in general is: try to structure your code such that they naturally occur as easily testable units.

    - For e.g. I had to pass a writer to functions _parseAddArgs_ and _parseQuoteArgs_, even though this is not strictly a functional style (side-effects to writer) because:
        1. Wanted my parsing logic to be easily testable
        2. The flags package _fs_ handles errors by directly writing to writer
    Know what you are doing and why you're doing it

    - All sub-command functions apart from these two are purely functional which made testing quite simple


- I had also tried to structure _add_ sub-command's core function in a non-functional manner (_addQuoteToStorage_ function, refer to git history). This resulted in convoluted tests and a weird mock structure + interface to test it out.

- Thought: Do you need to test for actually running the command using os.Exec()? Since you are mocking stderr and stdout using interfaces + testing all functional components, I don't think there is much need to do this.

- Haven't added tests for the _driver_ functions (_HandleAdd_, _HandleQuote_, _handleCmd_). Got a bit lazy + I feel I've got all I could out of this exercise

## Nit

- usage message for parent command will always name it as _quoter_ even if compiled binary is named differentky - might confuse user




