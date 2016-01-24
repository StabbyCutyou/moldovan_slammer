Moldovan Slammer
==================

The Moldovan Slammer is comprised of two packages: Moldova, and the Slammer.

# Moldova
Moldova is a lightweight template interpreter, used to generate random values to plug into a template, as defined by a series of custom tokens.

It understands the tokens as defined further down in the document.

# Slammer
Slammer is a simple utility using Moldova, for load testing a database. You can give it a template SQL query, most likely an INSERT
statement, and use that to generate a large volume of traffic, each request having a different set of values placed into it. In this way,
the Slammer makes an excellent tool for massively loading fake data into a database for load testing.

# Notice

Experimental - could change at a notice. Or, without notice. Have fun!

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
MS_PAUSEINTERVAL=200us MS_ITERATIONS=200000 MS_CONNSTRING="root@tcp(127.0.0.1:3306)/my_db" MS_INPUT="INSERT INTO floof VALUES ('{guid}','{guid:0}','{country}',{int:-2000:0},{int:100:1000},{float:-1000.0:-540.0},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL,-3)" ./moldovan_slammer
```

This would provide sample output like the following:

```sql
INSERT INTO floof VALUES ('03ad6a7b-a09a-4ede-b410-7a07dc868d0c','03ad6a7b-a09a-4ede-b410-7a07dc868d0c','BI',-1173,717,-1185.063842,32,'2016-01-23T14:50:43-05:00','2016-01-23T14:50:43-05:00','DS',NULL,-3)
INSERT INTO floof VALUES ('012ba1fa-38dd-4529-9d50-39eb59a3b495','012ba1fa-38dd-4529-9d50-39eb59a3b495','MX',-1582,555,-1259.542916,16,'2016-01-23T14:50:45-05:00','2016-01-23T14:50:45-05:00','KR',NULL,-3)
INSERT INTO floof VALUES ('188058a2-47d6-4cbc-93dc-b61cd3e1d29c','188058a2-47d6-4cbc-93dc-b61cd3e1d29c','FO',-1635,717,-1192.019471,34,'2016-01-23T14:50:47-05:00','2016-01-23T14:50:47-05:00','ER',NULL,-3)
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

I'll continue to add support for more random value categories, such as a general {time} field.

I also want to come up with a better internal design for how the interpreter is organized and architected, but I'm waiting until I have a richer feature set before I tackly an overall re-design of the current implementation. This likely won't happen until I split the libraries into Moldova and Slammer.

I also plan on making it possible to swap the current driver (mysql compatible databases only) for another one, longer term.

# License

Apache v2 - See LICENSE
