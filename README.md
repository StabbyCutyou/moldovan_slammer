Moldovan Slammer
==================

The Moldovan Slammer is used to generate random values to plug into a sample input string, which is meant to be something you'd run against a database (like bulk loading insert statements full of random data, hence the name "slammer"). It understands the tokens as defined below:

# Notice

Experimental - could change at a notice. Have fun!

# Tokens

## {guid:ordinal}

Slammer will replace any instance of {guid} with a guid, by shelling out to `uuidgen`, available on linux and osx. This is because Golang currently lacks
a built in standardized way to generate uuids. I'm willing to implement one
of the many third party options, once someone tells me which one is trustworthy
and a need is demonstrated

If you provide the :ordinal option, for the current line of text being generated,
you can have the Slammer insert an existing value, rather than a new one. For
example:

"{guid} - {guid:0}"

In this example, both guids will be replaced with the same value. This is a way
to back-reference existing generated values, for when you need something repeated.

## {now:ordinal}

Slammer will replace any instance of {now} with a string representation of Golangs
time.Now() function, formatted per the time.RFC3339 format.

{now} also supports the same :ordinal option as {guid}

## {integer:lower:upper}

Slammer will replace any instance of {integer} with a random int value, optionally between the range provided. The defaults, if not provided, are 0 to 100.

{integer} currently does not support :ordinal

## {float:lower:upper}

Slammer will replace any instance of {float} with a random Float64, optionally between the range provided. The defaults, if not provided, are 0.0 to 100.0

{float} currently does not support :ordinal

# {char:number:case}

Slammer will replace any instance of {char} with a randomly generated set of characters,
optionally up to the number specified by :number. The default value is 2.

{char} also takes the :case argument, which is either 'up' or 'down', like so

{char:5:up}
{char:2:down}

{char} currently does not support :ordinal, nor a mixing of cases
