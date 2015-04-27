# reactor

reactor is a helper Go package to generate a single page HTML app, using React. It allows the user to add javascript and css folders, which are (non-recursively) scanned and included in the ``<head>`` upon each request.

## Usage

```go
  app := reactor.New("FireApp")
  app.MapJavascriptFolder("public/js", "js")
  app.MapCssFolder("public/css", "css")

  io.WriteString(response, app.Html().String())
```

For a more complete example, check https://github.com/sparkymat/fire
