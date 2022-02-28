# shellspy

`shellspy` is a Go library and command-line client for the recording of sessions.

## Installing the command-line client

To install the client binary, run:

```
go get -u github.com/redscaresu/shellspy
```

## Using the command-line client

To use shellspy run:

```
> shellspy --port 9999
shellspy is running remotely on port 9999
```

OR

```
> go run cmd/main.go --mode local
shellspy is running locally
```