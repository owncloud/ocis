---
title: "Efficient Stat Polling"
date: 2020-03-03T10:31:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/architecture
geekdocFilePath: efficient-stat-polling.md
---

The fallback sync mechanism uses the ETag to determine which part of a sync tree needs to be checked by recursively descending into folders whose ETag has changed. The ETag can be calculated using a `stat()` call in the filesystem and we are going to explore how many `stat()` calls are necessary and how the number might be reduced.

## ETag propagation

What does ETag propagation mean? Whenever a file changes its content or metadata the ETag or "entity tag" changes. In the early days of ownCloud it was decided to extend this behavior to folders as well, which is outside of any WebDAV RFC specification. Nevertheless, here we are, using the ETag to reflect changes, not only on WebDAV resources but also WebDAV collections. The server will propagate the ETag change up to the root of the tree.

{{<mermaid class="text-center">}}
graph TD
  linkStyle default interpolate basis
  
  subgraph final ETag propagation
    ert3(( etag:N )) --- el3(( etag:O )) & er3(( etag:N ))
    er3 --- erl3(( etag:O )) & err3(( etag:N ))
  end

  subgraph first ETag propagation
    ert2(( etag:O )) --- el2(( etag:O )) & er2(( etag:N ))
    er2 --- erl2(( etag:O )) & err2(( etag:N ))
  end

  subgraph initial file change 
    ert(( etag:O )) --- el(( etag:O )) & er(( etag:O ))
    er --- erl(( etag:O )) & err(( etag:N ))
  end
{{</mermaid>}}

The old `etag:O` is replaced by propagating the new `etag:N` up to the root, where the client will pick it up and explore the tree by comparing the old ETags known to him with the state of the current ETags on the server. This form of sync is called *state based sync*.

## Single user sync
To let the client detect changes in the drive (a tree of files and folders) of a user, we rely on the ETag of every node in the tree. The discovery phase starts at the root of the tree and checks if the ETag has changed since the last discovery:
- if it is still the same nothing has changed inside the tree
- if it changed the client will compare the ETag of all immediate children and recursively descend into every node that changed

This works, because the server side will propagate ETag changes in the tree up to the root.

{{<mermaid class="text-center">}}
graph TD
  linkStyle default interpolate basis

  ec( client ) -->|"stat()"|ert

  subgraph  
    ert(( )) --- el(( )) & er(( ))
    er --- erl(( )) & err(( ))  
  end
{{</mermaid>}}

## Multiple users
On an ocis server there is not one user but many. Each of them may have one or more clients running. In the worst case all of them polling the ETag of his home root node every 30 seconds.

Keep in mind that etags are only propagated inside each distinct tree. No sharing is considered yet.

{{<mermaid class="text-center">}}
graph TD
  linkStyle default interpolate basis

  ec( client ) -->|"stat()"|ert

  subgraph  
    ert(( )) --- el(( )) & er(( ))
    er --- erl(( )) & err(( ))
  end
  
  mc( client ) -->|"stat()"|mrt

  subgraph  
    mrt(( )) --- ml(( )) & mr(( ))
    mr --- mrl(( )) & mrr(( ))
  end

  fc( client ) -->|"stat()"|frt

  subgraph  
    frt(( )) --- fl(( )) & fr(( ))
    fr --- frl(( )) & frr(( ))
  end
{{</mermaid>}}

## Sharing
*Storage providers* are responsible for persisting shares as close to the storage as possible.

One implementation may persist shares using ACLs, another might use custom extended attributes. The chosen implementation is storage specific and always a tradeoff between various requirements. Yet, the goal is to treat the storage provider as the single source of truth for all metadata. 

If users can bypass the storage provider using eg. `ssh` additional mechanisms needs to make sure no inconsistencies arise:
- the ETag must still be propagated in a tree, eg using inotify, a policy engine or workflows triggered by other means
- deleted files should land in the trash (eg. `rm` could be wrapped to move files to trash)
- overwriting files should create a new version ... other than a fuse fs I see no way of providing this for normal posix filesystems. Other storage backends that use the s3 protocol might provide versions natively.

The storage provider is also responsible for keeps track of references eg. using a shadow tree that users normally cannot see or representing them as symbolic links in the filesystem (Beware of symbolic link cycles. The clients are currently unaware of them and would flood the filesystem).

To prevent write amplification ETags must not propagate across references. When a file that was shared by einstein changes the ETag must not be propagated into any share recipients tree.

{{<mermaid class="text-center">}}
graph TD
  linkStyle default interpolate basis


  ec( einsteins client ) -->|"stat()"|ert

  subgraph  
    ml --- mlr(( ))
    mrt(( )) --- ml(( )) & mr(( ))
    mr --- mrl(( )) & mrr(( ))
  end

  mlr -. reference .-> er

  subgraph  
    ert(( )) --- el(( )) & er(( ))
    er --- erl(( )) & err(( ))
  end
  
  mc( maries client ) -->|"stat()"|mrt

{{</mermaid>}}

But how can maries client detect the change?

We are trading writes for reads: the client needs to stat the own tree & all shares or entry points into other storage trees.

It would require client changes that depend on the server side actually having an endpoint that can efficiently list all entry points into storages a user has access to including their current etag.

But having to list n storages might become a bottleneck anyway, so we are going to have the gateway calculate a virtual root ETag for all entry points a user has access to and cache that.

## Server Side Stat Polling
Every client polls the virtual root ETag (every 30 sec). The gateway will cache the virutal root ETag of every storage for 30 sec as well. That way every storage provider is only stated once every 30 sec (can be throttled dynamically to adapt to storage io load).


{{<mermaid class="text-center">}}
graph TD
  linkStyle default interpolate basis

  ec( client ) -->|"stat()"|evr

  subgraph gateway caching virtual etags
    evr(( ))
    mvr(( ))
    fvr(( ))
  end

  evr --- ert
  mvr --- mrt
  fvr --- frt

  subgraph  
    ert(( )) --- el(( )) & er(( ))
    er --- erl(( )) & err(( ))
  end
  
  mc( client ) -->|"stat()"|mvr

  subgraph  
    mrt(( )) --- ml(( )) & mr(( ))
    ml --- mlm(( ))
    mr --- mrl(( )) & mrr(( ))
  end

  mlm -.- er
  mvr -.- er

  fc( client ) -->|"stat()"|fvr

  subgraph  
    frt(( )) --- fl(( )) & fr(( ))
    fr --- frl(( )) & frr(( ))
  end

{{</mermaid>}}

Since the active clients will poll the etag for all active users the gateway will have their ETag cached. This is where sharing comes into play: The gateway also needs to stat the ETag of all other entry points ... or mount points. That may increases the number of stat like requests to storage providers by an order of magnitude.

### Ram considerations

For a single machine using a local posix storage the linux kernel already caches the inodes that contain the metadata that is necessary to calculate the ETag (even extended attributes are supported). With 4k inodes 256 nodes take 1Mb of RAM, 1k inodes take 4Mb and 1M inodes take 4Gb to completely cache the file metadata. For distributed filesystems a dedicated cache might make sense to prevent hammering it with stat like requests to calculate ETags.

### Bandwith considerations

The bandwith for a single machine might be another bottleneck. Consider a propfind request with roughly 500 bytes and a response with roughly 800 bytes in size:
- At 100Mbit (~10Mb/s) you can receive 20 000 PROPFIND requests
- At 1000Mbit (~100Mb/s) you can receive 200 000 PROPFIND requests
- At 10Gbit (~1Gb/s) you can receive 2 000 000 PROPFIND requests

This can be scaled by adding more gateways and sharding users because these components are stateless.

## Share mount point polling cache
What can we do to reduce the number of stat calls to storage providers. Well, the gateway queries the share manager for all mounted shares of a user (or all entry points, not only the users own root/home). The share references contain the storage provider that contains the share. If every user has it's own storage provider id the gateway could check in its own cache if the storage root etag has changed. It will be up to date because another client likely already polled for its etag.
This would reduce the number of necessary stat requests to active storages.

### Active share node cache invalidation
We can extend the lifetime of share ETag cache entries and only invalidate them when the root of the storage that contains them changes its ETag. That would reduce the number of stat requests to the number of active users.

### Push notifications
We can further enhance this by sending push notifications when the root of a storage changes. Which is becoming increasingly necessary for mobile devices anyway.
