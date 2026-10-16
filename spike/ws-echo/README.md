# WebSocket echo idle test (spike, issue #15)

Throwaway. Remove or move into the validation ADR once it has served its purpose.

```
POKER_HTTP_ADDR=:8080 go run ./spike/ws-echo
```

Open `/` (optionally `/?interval=30`), leave it open, and watch the state, reconnect counter and log.
Try several intervals (for example 15, 30, 60, 120 s) to find where an intermediate proxy drops idle connections.
The server logs connection lifetime and close reason only, never addresses.
