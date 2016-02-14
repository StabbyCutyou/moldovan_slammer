Moldovan Slammer
==================

The Moldovan Slammer a helper / tutorial repo for using [Moldova](https://github.com/StabbyCutyou/moldova), and the [Slammer](https://github.com/StabbyCutyou/slammer).
It's designed to demonstrate how to use both of the libraries together to generate random data for database load testing.

If you're not familiar with those packages, i'd suggest checking them out first. But at a high level,
Moldova can be used to generate random data for insert or select statements, and Slammer accepts sql
statements as inputs, running in a worker pool and giving a report on the latency, throughput, and error
rate.

# Using them together

You can use them together in a few ways. The first is to simply pipe the output of
moldova into the slammer, like so

```bash
moldova -t "INSERT INTO floof VALUES ('{guid}','{guid:ordinal:0}','{country}',{int:min:-2000|max:0},{int:min:100:|max:1000},{float:min:-1000.0|max:-540.0},{int:min:1|max:40},'{now}','{now:ordinal:0}','{country:case:up}',NULL,-3)" -n 100 | slammer -c "root@tcp(10.248.5.220:3306)/tapjoy_db" -p 200us -w 2
```

This will generate a new list of random data for every insert statement

You could also pre-generate a series of statements, and issue them against slammer sepparately, like so

```bash
moldova -t "INSERT INTO floof VALUES ('{guid}','{guid:ordinal:0}','{country}',{int:min:-2000|max:0},{int:min:100:|max:1000},{float:min:-1000.0|max:-540.0},{int:min:1|max:40},'{now}','{now:ordinal:0}','{country:case:up}',NULL,-3)" -n 100 > dbdata
slammer -c "root@tcp(10.248.5.220:3306)/tapjoy_db" -p 200us -w 2 < dbdata
```

# License

Apache v2 - See LICENSE
