# Testing Vektor Servers ✅

`vtest` is a powerful package for testing your Vektor servers without running an HTTP server bound to a port. `vtest` uses the standard Go `testing` package, which lets you integrate server route testing into your test suite. Let's explore a simple wordcount API server and add some tests to it.

All of the code on this page is runnable and can be found on [GitHub](https://github.com/suborbital/vektor/tree/main/docs/examples/wordcount).

## The Server
Our API consists of helper functions, a middleware function to call the helper functions, and a single `POST` endpoint.

The `Wordcount` type is a simple wrapper around `string` that adds some handy helper methods similar to the ubiquitous UNIX program `wc`.
```go
// wordcount.go
type Wordcount string

func (w Wordcount) Words() int {
	return len(strings.Fields(string(w)))
}

func (w Wordcount) Lines() int {
	return len(strings.Split(string(w), "\n"))
}

func (w Wordcount) Characters() int {
	runes := []rune(w)
	return len(runes)
}
```

There are a few interesting things to point out in our server that are relevant to testing. The `setupServer()` function is reused later to setup our testing server in the same way it is used in `main()`. This is useful if you have a more complex routing setup.

The other thing to note here is the `WCResponse` struct. We can reuse it directly to make writing tests a breeze.

```go
// server.go
func setupServer() *vk.Server {
	server := vk.New(vk.UseAppName("wordcount"), vk.UseHTTPPort(9090))
	api := vk.Group("/api/v1").Before(createWordcountMiddleware)
	api.POST("/wc", handlePost)

	server.AddGroup(api)

	return server
}

func createWordcountMiddleware(r *http.Request, ctx *vk.Ctx) error {
	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	wc := Wordcount(string(text))
	ctx.Set(wordcountCtxKey, wc)

	return nil
}

type WCResponse struct {
	Words      int `json:"words"`
	Lines      int `json:"lines"`
	Characters int `json:"chars"`
}

func NewWCResponse(wc Wordcount) *WCResponse {
	return &WCResponse{
		Words:      wc.Words(),
		Lines:      wc.Lines(),
		Characters: wc.Characters(),
	}
}

func handlePost(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	wc := ctx.Get(wordcountCtxKey).(Wordcount)

	return vk.R(http.StatusOK, NewWCResponse(wc)), nil
}

func main() {
	server := setupServer()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
```

As a quick check, let's run our wordcount server.

```
cd docs/examples/wordcount
go run .
```

You should see the default logger output:
```
{"log_message":"(W) configured to use HTTP with no TLS","timestamp":"2021-08-16T10:45:46.273974908-04:00","level":2}
{"log_message":"(I) starting wordcount ...","timestamp":"2021-08-16T10:45:46.274132326-04:00","level":3}
{"log_message":"(I) serving on :9090","timestamp":"2021-08-16T10:45:46.274144389-04:00","level":3}
```

Let's send it a request. In a separate terminal run:
```
curl http://localhost:9090/api/v1/wc -d "Hello, word count"
```

Which should produce:
```json
{"words":3,"lines":1,"chars":17}
```

## Writing Tests
Now that our server seems to be running as expected, let's write tests with the `vtest` package.

The `vtest` package is part of Vektor and can be found in at [github.com/suborbital/vektor/tree/main/vtest](https://github.com/suborbital/vektor/tree/main/vtest).

Tests written with `vtest` use the usual Go testing idioms. If you're not familiar with the `testing` package, you can [read about it here](https://golang.org/doc/tutorial/add-a-test) as part of the official Go tutorial.  

Here is a complete test function for a Vektor server. Let's break it down.
```go
func TestWordcount(t *testing.T) {
	server := setupServer()

	vt := vtest.New(server)

	body := strings.NewReader("There's a starman waiting in the sky\nHe'd like to come and meet us")

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wc", body)

	vt.Do(req, t).
		AssertStatus(http.StatusOK).
		AssertJSON(WCResponse{
			Words:      14,
			Lines:      2,
			Characters: 66,
		})
}
```

The only thing different from creating a regular Vektor server is that we construct a `vtest.VTest` struct with `vtest.New()`.
```go
server := setupServer()
vt := vtest.New(server)
```

Next, we create a normal HTTP request with Go standard library functions. Nothing fancy here.
```go
body := strings.NewReader("Hello!")
req, err := http.NewRequest(http.MethodPost, "/api/v1/wc", body)
```

And finally, the interesting part, passing the request to our `vt` object. Note that we use the exact same struct here as was used in the server itself. This is super useful if you have defined a custom `MarshalJSON` method, for example. 
```go
vt.Do(req, t).
    AssertStatus(http.StatusOK).
    AssertJSON(WCResponse{
        Words:      1,
        Lines:      1,
        Characters: 6,
    })
```

`Do()` returns a `*Response` object, as do each of the `Assert` methods of `Response`, which lets us chain assertions together without rerunning the request.
```go
func (vt *VTest) Do(req *http.Request, t *testing.T) *Response
```

## Running Tests

Just as you would test any other standard Go package, simply run:

```
go test . -v
```

The `-v` (verbose) flag lets us see details of tests, even if they pass. One thing to note here is that certain assertion helpers, such as `AssertJSON()`, perform more than one test. `vtest` creates subtests automatically for you in these cases. These are shown indented below. Another example is the `AssertHeaders()` helper, which performs a nested subtest for each header.

`vtest` tries to stick to standard Go testing idioms while making common testing tasks easy to do.
```
=== RUN   TestWordcount
=== RUN   TestWordcount/status
=== RUN   TestWordcount/JSON
=== RUN   TestWordcount/JSON/Content-Type
=== RUN   TestWordcount/JSON/body
--- PASS: TestWordcount (0.00s)
    --- PASS: TestWordcount/status (0.00s)
    --- PASS: TestWordcount/JSON (0.00s)
        --- PASS: TestWordcount/JSON/Content-Type (0.00s)
        --- PASS: TestWordcount/JSON/body (0.00s)
PASS
ok  	github.com/suborbital/vektor/docs/examples/wordcount	0.003s
```

All of our tests have passed. Great! 

## Documentation
Further documentation for `vtest` can always be found in godoc. There are also more examples in the `vk/test` and `vtest/` directories. 

```
go doc github.com/suborbital/vektor/vtest Response
```
