This repository is a reproduction for [Chromium issue 1415291](https://bugs.chromium.org/p/chromium/issues/detail?id=1415291).
The issue is that since Chromium disabled HTTP2 PUSH there is no way to preload content negotiated data (i.e., where data
depends on the `Accept` header) which is a problem when you have a server which serves both HTML and JSON at the same path
and you want to preload JSON when HTML is requested.

To run, you should use [Go](https://go.dev/):

```
go run main.go
```

And then open [http://localhost:8000/](http://localhost:8000/) in Chromium. Open network tab in DevTools.

## Expected

You should see three requests being made very soon after the other, one for `/` for `text/html`, another for
`/` for `application/json`, and another for `/data.json` for `application/json`. After 2 seconds that JSON data
should be shown in the page when JavaScript calls `fetch` and gets preloaded data without doing another
server equest.

## Actual

You see three requests, one for `/`, anther for `/` and the third for `/data.json`, but all of them
`text/html`. After 2 seconds no additional requests are made (good), but `fetch` gets HTML responses
which are invalid JSON so error is shown in the page by JavaScript (bad).

## Discussion

There are multiple issues here:

- Preload requests are made with `Accept` set to `*/*` so HTML responses are returned. Instead,
  `type="application/json"` from `Link` header should be respected to load JSON responses.
- When JavaScript calls `fetch` 2 seconds later, it is called with `Accept` header set to
  `application/json` but Chromium still returns preloaded HTML content instead of doing another
  request with `application/json`. So not just that performance is degraded (invalid and unnecessary
  preloading of HTML is done) but also correctness of execution of the page is impacted.
- No preloading happens if just `</>; rel="preload"; as="fetch"; type="application/json"; crossorigin="anonymous"`
  `Link` header is issued by the server. If both `</>; rel="preload"; as="fetch"; type="application/json"; crossorigin="anonymous"`
  and `</data.json>; rel="preload"; as="fetch"; type="application/json"; crossorigin="anonymous"` headers are issued, then
  both of them are preloaded. (This is the reason why I included additional `/data.json` loading. Originally the plan was
  to have loading of HTML and JSON just from `/`.)

Tested on Chromium 114.0.5735.198.

Firefox 115.0.2 behaves the same but it does not require two `Link` headers to load from `/` and
it still has HTTP2 PUSH.
