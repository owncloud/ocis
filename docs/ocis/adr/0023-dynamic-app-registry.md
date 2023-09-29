* Status: proposed
* Deciders: @butonic, @dragonchaser
* Date: 2023-09-28

## Context and Problem Statement

The intent of the CS3 API is to link enterprise file share and sync platforms with storage and application providers. The CS3 [RegistryAPI](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.RegistryAPI) allows registering CS3 [AppProviders](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.ProviderInfo). However, the API does not explain if and how often an AppProvider needs to reregister. In K8s ip addresses may change anytime and pods can get started and removed at any time. Without a mechanism to reupdate the App Registry when that happens the metadata of an app provider - the filetype it can handle and the endpoint to use - will be invalid or forgotten.

## Out of scope
Additional app properties like web UI elements or scripts that Web could use to dynamically load extensions. That will be part of a seperate ADR.

## Decision Drivers

## Considered Options

* [Dynamically reregister App Providers](#dynamic-registry)
* [Switch to a Go Micro based Registry](#micro-registry)

## Decision Outcome

Chosen option: [Dynamically reregister App Providers](#dynamic-registry)

### Positive Consequences

* App-providers can register service endpoints, supported mimetypes
* App-providers gain the ability to refresh their data in the registry, the registry gains a ttl for each service resulting in automatic removal & cleanup of dead services.
* CS3 API usage clarified and we are not bypassing by reading the app provider metadata from the go micro registry.

### Negative Consequences

* We need to update the CS3 API spec and add a TTL to the [ProviderInfo](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.ProviderInfo) so it can can be used in the [ListAppProvidersResponse](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.ListAppProvidersResponse) but more importantly in the [AddAppProviderRequest](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.AddAppProviderRequest) so App providers can optionally request a TTL. Finally, the [AddAppProviderResponse](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.AddAppProviderResponse) also needs to return a TTL so the App Provider knows when to reregister.
* The static app registry implementation needs a config option for a TTL

## Pros and Cons of the Options

### Dynamically reregister App Providers
The App Registry would use a TTL to determine when to forget App Providers and App Providers would have to refresh their registration periodically.

* Good, because admin has full control over the default applications which are configured in the app registry.
* Good, because App Providers can dynamically update the extensions / mimetypes they can handle without ocis having to restart, solving the impracticability in k8s scenarios mentioned above.
  [jfd: I don't know if a 60sec time window is good enough ... or if we can make the app providers listen for changes in the service registry so that they will only update the registry if it gets restarted. So instead of a timer in the app provider and a ttl in the registry set ttl to infinity and use a registry watcher to determine when to update. Or do both: ttl and timer as a daily fallback and a registry to immediately trigger reregistering]

### Switch to a Go Micro based Registry
The idea is to use an existing go micro service registry implementation and transport the mimetypes the approvider supports via m̀icro service metadata.

* Good, because it allows app-providers to register all necessarry data like service endpoints, supported mimetypes.
* Good, because it allows reregistering all necessary data upon service restart, which solves the impracticability in k8s scenarios mentioned above.
* Bad because it obsoletes the existing CS3 API [AddAppProviderRequest](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.AddAppProviderRequest) and would require fallback handling in case an AppProvider does not use go micro. This goes against the idea of the CS3 API to link enterprise file share and sync platforms with storage and application providers.

## Links
* https://github.com/owncloud/ocis/issues/3832
