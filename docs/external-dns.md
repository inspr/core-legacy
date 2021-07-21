

### steps to install the external dns in the cluster


```
helm install my-release \
--set provider=google \
--set clusterdomain=inspr.dev \
--set policy=sync \
--set google.project=insprlabs \
bitnami/external-dns
```
