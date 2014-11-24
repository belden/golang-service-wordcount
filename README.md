## INSTALLATION

### Install Dependencies

Well, I've never built anything in Go before, so I searched around for a REST API for Go that
looked sensible to me. I found [sleepy](http://dougblack.io/words/a-restful-micro-framework-in-go.html),
which is a very small framework - about 100 lines of Go.

I thought I'd go with [sleepy](github.com/dougblack/sleepy), but ended up using plain old `net/http`. The code is quite messy as a result - dispatching of different request methods is ugly. I wanted to code towards functionality more than towards beauty.

### Run the test suite

Test that you can run the test suite:

```
perl -MJSON -le 'print "found JSON.pm"'
```

If the above command runs successfully, then your local version of `perl` has a `JSON.pm` installed. If the above fails, then you may install `JSON.pm` in any number of ways - the easiest is to run

```
cpan JSON
```

Now you can fire up the tests. The test suite starts the service, and makes various assertions about interacting with the service, then shuts the service down pretty roughly.

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
curl -s -F file=@/tmp/test-file.txt -F filename=/tmp/test-file.txt http://localhost:3000/files
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

#### GET /wordcounts

Returns the list of filenames in cache.

Sample request:

```
curl http://localhost:3000/files
```

Sample response:

```
["foo","bar","/etc/hosts"]
```

#### GET /wordcounts?filename=:filename

Return the cached data for the file `:filename`.

Sample request:

```
curl http://localhost:3000/files?filename=/tmp/test-file.txt
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

#### DELETE /wordcounts?filename=:filename

Remove `:filename` from the cache.

Sample request:

```
curl -X DELETE http://localhost:3000/files?filename=/tmp/test-file.txt
```

Sample response:

```
null
```

(I'm not happy with this return value.)

### Error handling

Close to none. Requests that cause errors:

```
curl -X DELETE http://localhost:3000/files    # no filename given
curl -X PUT http://localhost:3000/files       # PUT will fall through POST and fail
```

### Logging

Close to none. Certianly nothing standard. This doesn't use standard log formats for inbound requests.

### Unit Tests

None yet - the Golang testing tools look pretty neat, but I haven't done anything here yet.
