## INSTALLATION

### Install Dependencies

Well, I've never built anything in Go before, so I searched around for a REST API for Go that
looked sensible to me. I found [sleepy](http://dougblack.io/words/a-restful-micro-framework-in-go.html),
which is a very small framework - about 100 lines of Go. I don't know if this is particularly idiomatic
Go, but it's what I went with.

```
go get github.com/dougblack/sleepy
```

### Run the Service

```
go run service-wordcount.go
```

The service starts on port 3000.

### Resources

#### POST /wordcounts

Reads the uploaded file `filename` and returns a JSON structure of the assigned filename, the
count of words, and the frequency of individual words.

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
  "filename": "/tmp/test-file.txt",
  "wordcount": 15,
  "words": {
    "Hello": 1,
    "world": 1,
    "This": 2,
    "is": 1,
    "a": 1,
    "sentence": 1,
    "with": 1,
    "extra": 1,
    "spaces": 1,
    "Is": 1,
    "this": 1,
    "the": 1,
    "same": 1,
    "as": 1,
  }
}
```
