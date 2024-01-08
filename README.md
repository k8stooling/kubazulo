<img src="https://github.com/maikbrauer/kubazulo/assets/53018978/9a6f35a3-e233-4b07-93cd-76bc560489b4" width="64" />

# kubazulo
Kubeconfig Authentication Helper for Kubernetes API-Server in cunjunction with kubectl

## Description
kubazulo is a client-go credential (exec) plugin implementing azure authentication. It plugs in seemless into the process of communicating to the kubernetes API-Server.

For this the kubeconfig needs to be adapted.

## Setup the k8s OIDC Provider 

kubazulo can be used to authenticate to general kubernetes clusters using Azure Active Directory as an OIDC provider.

1. Create an AAD Enterprise Application and the corresponding App Registration. Check the Allow public client flows checkbox. Configure groups to be included in the response. Take a note of the directory (tenant) ID as $AAD_TENANT_ID and the application (client) ID as $AAD_CLIENT_ID

2. Configure the API server with the following flags:

* Issuer URL: --oidc-issuer-url=https://sts.windows.net/$AAD_TENANT_ID/
* Client ID: --oidc-client-id=$AAD_CLIENT_ID
* Username claim: --oidc-username-claim=upn
* Group claim --oidc-groups-claim=groups

>See the [kubernetes](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server) docs for optional flags.

3. Configure the Exec plugin with kubelogin to use the application from the first step:


### Configure for Standalone Flow (Default)
```
kubectl config set-credentials "kubazulo-azuread" \
  --exec-api-version=client.authentication.k8s.io/v1 \
  --exec-command=kubazulo \
  --exec-arg=get-token \
  --exec-arg=--client-id \
  --exec-arg=$AAD_CLIENT_ID \
  --exec-arg=--tenant-id \
  --exec-arg=$AAD_TENANT_ID
```

### Configure for Intermediate Flow (Advanced)
```
kubectl config set-credentials "kubazulo-azuread" \
  --exec-api-version=client.authentication.k8s.io/v1beta1 \
  --exec-command=kubazulo \
  --exec-arg=get-token \
  --exec-arg=--client-id \
  --exec-arg=$AAD_CLIENT_ID \
  --exec-arg=--tenant-id \
  --exec-arg=$AAD_TENANT_ID \
  --exec-arg=--loginmode \
  --exec-arg=interactive \
  --exec-arg=--intermediate \
  --exec-arg=true \
  --exec-arg=--api-token-endpoint \
  --exec-arg=$APIGW_ENDPOINT
```
>Please DON'T FORGET TO SET THE OS-Environment Variables |


4. Use this credential to connect to the cluster:

## Command Argument

* get-token
* version

## Command Flags

| Parameter            | Description                              | Mandatory | Default     |
|----------------------|------------------------------------------|:---------:|-------------|
| --client-id          | Azure Application-ID                     |:heavy_check_mark: | n/a         |
| --tenant-id          | Azure Tenant-ID                          |:heavy_check_mark: | n/a         |
| --force-login        | Re-Usage of Browser Session data         |:x: | false       |
| --loopbackport       | Customize local callback listener        |:x: | 58433       |
| --loginmode          | Set the Authentication Flow mode         |:x: | interactive |
| --intermediate       | Activate another Token fetcher Endpoint  |:x: | false       |
| --api-token-endpoint | Define Endpoint from where it gets Token |:x: | n/a         |

## Logging

kubazulo will also log the operations it is doing to the following folder
`$HOME/.kube/kubazulo/application.log`

## References
### kubectl Command Installation
https://kubernetes.io/docs/tasks/tools/

```
kubectl config set-context "$CLUSTER_NAME" --cluster="$CLUSTER_NAME" --user=kubazulo-azuread
kubectl config use-context "$CLUSTER_NAME"
```
