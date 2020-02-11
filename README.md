Demo for using Keycloak to secure a go backend (using jwt)
===

## Whats in the box?
Simple showcase how to use Keycloak in combination with a go app to utilize JWT to secure your api


## How to run it?
```bash
KEYCLOAK_URL=https://keycloak.example.com/auth/realms/master KEYCLOAK_CLIENT=your-client-id go run *.go
```

You will need a bearer token to send to your api with `curl` like this:

```bash
curl -H 'Authorization: Bearer <token from keycloak>' http://localhost:3000/hello
```


## How does it work?
This demo is relying on [go-idc](https://github.com/coreos/go-oidc) for the token validation etc; this logic is used
inside a middleware to validate the token and given response according to the result of the validations.

It also enhances the request context with a few details from the jwt token, that are set by keycloak.


## Disclaimer
**This is just a quick and dirty demo. There is no real permission check implemented. It simply checks if the user has
 access to the client and validates the token is signed by your keycloak instance. 
 The errors are passed to the user as is, so not that super secure.**
 
*Keep that in mind when looking at this quick and dirty demo and dont copy/paste it like a fool!*