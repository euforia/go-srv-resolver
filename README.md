# go-srv-resolver
This package implements a resolver to perform SRV record lookups.  It can be used to communicate with
services uses service discovery.


### Example

```

import "github.com/euforia/go-srv-resolver"

// Connect to consul
rslv := resolver.NewResolver(8600, "127.0.0.1")

// Perform SRV lookup.
resp, err:= rslv.Lookup("consul.service.consul")
...

for _, r := range resp {
    fmt.Printf("host=%s ip=%s port=%d\n", r.Hostname, r.IP, r.Port)
}
```
