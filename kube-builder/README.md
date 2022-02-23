# Kube-Builder

- `Process`
  - `Manager`
      - `Client`
      - `Cache`
          - `Controller`
              - `Predicate` : filter a stream of events. passing only those that require action to the reconciler
               - `Reconciler`
       - `Webhook`
         - `Defaulter`
         - `Validator`