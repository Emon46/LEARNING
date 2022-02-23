#Resources

https://helm.sh/docs/intro/using_helm/#more-installation-methods

# Concept:
- `Chart:` A Chart is a Helm package. It contains all of the resource definitions necessary to run an application, tool, or service inside of a Kubernetes cluster.
- `Repository:` A Repository is the place where charts can be collected and shared.
- `Release:` A Release is an instance of a chart running in a Kubernetes cluster. One chart can often be installed many times into the same cluster. And each time it is installed, a new release is created. 

```yaml
 helm create <chart-name>
```
will create helm chart 


```yaml
helm template <release-name> <chart-name>
```
```yaml
 helm install [NAME] [CHART] [flags]
```