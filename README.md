Moldovan Slammer
==================

The Moldovan Slammer is used to generate random values to plug into a sample input string, which is meant to be something you'd run against a database (like bulk loading insert statements full of random data, hence the name "slammer"). It understands the tokens as defined below:

# Notice

Experimental - could change at a notice. Have fun!

# Example usage

Right now, the Slammer gets configured via environment variables. This may change to a flag-based approach in the future.

### MS_CONNSTRING
The connection string for your database

### MS_INPUT
The template you wish to turn into a series of inserts - see below for examples

### MS_PAUSEINTERVAL
A valid time.Duration parsable string representing how long to wait in between each run of the statement.

### MS_ITERATIONS
How many times to run the statements overall

## Example

```bash
MS_PAUSEINTERVAL=200us MS_ITERATIONS=200000 MS_CONNSTRING="root@tcp(127.0.0.1:3306)/my_db" MS_INPUT="INSERT INTO floof VALUES ('{guid}','{guid:0}','{country}',{int:-2000:0},{int:100:1000},{int:100:1000},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL,-3)" ./moldovan_slammer
```

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

## {char:number:case}

Slammer will replace any instance of {char} with a randomly generated set of characters,
optionally up to the number specified by :number. The default value is 2.

{char} also takes the :case argument, which is either 'up' or 'down', like so

{char:5:up}
{char:2:down}

{char} currently does not support :ordinal, nor a mixing of cases

## {country:case:ordinal}

Slammer will replace any instance of {country} with an ISO 3166-1 alpha-2 country code.

{country} supports the same :case argument as {char}. The default value is "up"

{country} also supports the same :ordinal argument as {guid}. Because of how the template is interpreted, you must provide the optional :case argument if you are also to specify an ordinal.

# Roadmap

Eventually I will likely break this apart into 2 packages, or at the very least 2 sub-packages: Moldova, the light weight template interpreter, and Slammer, the SQL data-loader.

I'll also continue to add support for more random value categories, such as a general {time} field, as well as add support for generating data using the full unicode table, and not just the ASCII characters a - z.

I also want to come up with a better internal design for how the interpreter is organized and architected, but I'm waiting until I have a richer feature set before I tackly an overall re-design of the current implementation. This likely won't happen until I split the libraries into Moldova and Slammer.
