# Wish and Bubble Tea hello-world (spike, issue #13)

Throwaway. Remove or move into the validation ADR once it has served its purpose.

```
POKER_SSH_ADDR=127.0.0.1:2222 go run ./spike/ssh-hello
ssh -t -p 2222 localhost
```

Shows the window size and a resize counter; resize the terminal to see it update, `q` quits.

The host key is read from `POKER_SSH_HOST_KEY_PATH` (generated there if the file is missing).
Without it an ephemeral key is generated with a warning, so clients will see a changed host key on every restart.
No key material lives in the repo. The server logs nothing about clients.
