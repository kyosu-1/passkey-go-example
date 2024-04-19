# passkey-go-example
Try passkey in go.

Note that this implementation uses the following two libraries:

- https://github.com/go-webauthn/webauthn
- https://github.com/MasterKale/SimpleWebAuthn

## How to run

```bash
go run cmd/server/main.go
```

access to `http://localhost:8080/` and try to register and login.

## Reference

- https://www.w3.org/TR/webauthn-3/
