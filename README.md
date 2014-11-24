## INSTALLATION

### Install Dependencies

Well, I've never built anything in Go before, so I searched around for a REST API for Go that
looked sensible to me. I found [sleepy](http://dougblack.io/words/a-restful-micro-framework-in-go.html),
which is a very small framework - about 100 lines of Go.

```
go get github.com/dougblack/sleepy
```

### Run the test suite

The test suite starts the service, and makes various assertions about interacting with the service.

```
./test.sh
```

### Run the Service

```
go run service-wordcount.go --port 8080
```

The service starts on port 3000 by default.

### Resources

#### POST /wordcounts

Reads the uploaded file `filename` and returns a JSON structure of the assigned filename, the
count of words, and the frequency of individual words. Words are normalized to lower-case for
the purpose of frequency checking.

Sample request:
```
# make a test file
cat <<'END' > /tmp/test-file.txt
Hello world!
This is  a sentence with extra     spaces.
Is  "this" the same as "This"?
END

# send it to the running service
curl -s -F file=@/tmp/test-file.txt -F filename=/tmp/test-file.txt http://localhost:3000/
```

Sample response:

```
{
  "Total": 15,
  "Words": {
    "a"        : 1,
    "as"       : 1,
    "extra"    : 1,
    "hello"    : 1,
    "is"       : 2,
    "same"     : 1,
    "sentence" : 1,
    "spaces"   : 1,
    "the"      : 1,
    "this"     : 3,
    "with"     : 1,
    "world"    : 1,
  }
}
```
