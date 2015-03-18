# negroni-rollbar
A negroni middleware for rollbar

The middleware forwards all panics to rollbar.com.

```
import "github.com/jfbus/negroni-rollbar"

func main() {
	n := negroni.Classic()
    n.Use(rollbar.Report(rollbar.Config{Token: ROLLBAR_TOKEN}))
}
```

rollbar.Report recovers panics, the default Recovery handler does nothing.
