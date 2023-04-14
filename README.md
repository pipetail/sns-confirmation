# SNS confirmation

Confirmation requests can be easily distinguished by the Header so confirmation requests can be routed to the dedicated service.

```
X-Amz-Sns-Message-Type: SubscriptionConfirmation
```

## Traefik example

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: main
spec:
  entryPoints:
  - web
  - websecure
  routes:
  - kind: Rule
    match: Host(`some.site.com.uk) && PathPrefix(`/api/v1/commands`) && (HeadersRegexp(`X-Amz-Sns-Message-Type`, `SubscriptionConfirmation`))
    middlewares: []
    services:
    - name: sns-n-confirmation
      port: 80
```
