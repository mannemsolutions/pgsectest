# pgsectest
A tool to run security checks against postgres and return a score

## The origin
We wanted to run security tests, but ajutomated across all instances we manage, and deliver results to our clients.
And thus [pgsectest](https://github.com/MannemSolutions/pgsectest) was born.

## Downloading pgsectest
The most straight forward way is to download [pgsectest](https://github.com/MannemSolutions/pgsectest) directly from the [github release page](https://github.com/MannemSolutions/pgsectest/releases).
But there are other options, like
- using the [container image from dockerhub](https://hub.docker.com/repository/docker/mannemsolutions/pgsectest/general)
- direct build from source (if you feel you must)

Please refer to [our download instructions](DOWNLOAD_AND_RUN.md) for more details on all options.

## Usage
After downloading the binary to a folder in your path, you can run pgsectest with a command like:
```bash
pgsectest ./mytest*.yml ./andonemoretest.yml
```
Or using stdin:
```bash
cat ./mytests*.yml | pgsectest
```

## Verbosity
You can improve verbosity of output by adding one or more -v arguments:

```bash
pgsectest ./mytest*.yml ./andonemoretest.yml -vvv
```
Number of V's | Output
--- | ---
0 | Only end score
1 | Also score for failed tests
2 | Also advisory and url for failed tests
3 | Also score for succeded tests (max score)

## Defining your tests
A more detailed description can be found in [our test definition guide](TESTS.md).

TLDR; you can define one or more test chapters as yaml documents (separated by the '---' yaml doc separator).
Each test chapter can have the following information defined:
- a dsn, whith all connection details to connect to postgres.
  - **Note** that instead of configuring in this chapter, the [libpq environment variables](https://www.postgresql.org/docs/current/libpq-envars.html) can also be used, but options configured in this chapter take precedence.
- You can set the number of retries, delay and debugging options
- Each test can define
  - a name (defaults to the query when not set),
  - a query for the dividend and a query for the divisor
  - an advisory how to improve your score
  - a url for more information
  - the expected result (a list of key/value pairs)
  - the option to reverse the outcome (Ok results are counted as errors and vice versa)

Some example test definitions can be found in the [testdata](./testdata/) folder.
