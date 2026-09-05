# fulfillment

A Helm chart for fulfillment

**Homepage:** <https://github.com/giantswarm/fulfillment>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| image.registry | string | `"gsoci.azurecr.io"` |  |
| image.name | string | `"giantswarm/fulfillment"` |  |
| image.tag | string | `""` |  |
| aws.access_key_id | string | `""` |  |
| aws.secret_access_key | string | `""` |  |
| ingress.enabled | bool | `true` |  |
| ingress.annotations."cert-manager.io/cluster-issuer" | string | `"letsencrypt-giantswarm"` |  |
| ingress.annotations."kubernetes.io/tls-acme" | string | `"true"` |  |
| slack.token | string | `""` |  |
| route.enabled | bool | `false` |  |
| route.name | string | `"fulfillment"` |  |
| route.kind | string | `"HTTPRoute"` |  |
| route.annotations | object | `{}` |  |
| route.labels | object | `{}` |  |
| route.hostnames | list | `[]` |  |
| route.parentRefs | list | `[]` |  |
| route.matches[0].path.type | string | `"PathPrefix"` |  |
| route.matches[0].path.value | string | `"/"` |  |
| route.filters | list | `[]` |  |
| route.additionalRules | list | `[]` |  |
| route.securityPolicy.enabled | bool | `false` |  |
| route.securityPolicy.labels | object | `{}` |  |
| route.securityPolicy.annotations | object | `{}` |  |
