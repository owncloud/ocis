---
title: Port Ranges
date: 2018-05-02T00:00:00+00:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions
geekdocFilePath: port-ranges.go
geekdocCollapseSection: true
---

oCIS extensions often need a port to expose their services to other services or the outside world.
As users may have many different extensions running on the same machine, we should track port usage in the oCIS ecosystem. In the best case we ensure that each extension uses a non colliding port range, to make life of users easier.

This page tracks the knowingly used port ranges.

Feel free to "reserve" a free port range when you're developing an extension by adding your extension to the list (see the edit button in the top right corner).

If you're developing a non-public extension, we recommend using ports outside of the ranges listed below.

We also suggest to use the last port in your extensions' range as a debug/metrics port.


## Allocations

| Port range | Service                                                                       |
| ---------- | ----------------------------------------------------------------------------- |
| 9000-9010  | [ocis](https://github.com/owncloud/ocis)                                      |
| 9100-9104  | [web](https://github.com/owncloud/ocis/tree/master/web)                       |
| 9105-9109  | [hello](https://github.com/owncloud/ocis-hello)                               |
| 9110-9114  | [ocs](https://github.com/owncloud/ocis/tree/master/ocs)                       |
| 9115-9119  | [webdav](https://github.com/owncloud/ocis/tree/master/webdav)                 |
| 9120-9124  | [graph](https://github.com/owncloud/ocis/tree/master/graph)                   |
| 9125-9129  | FREE (formerly used by glauth)                                                |
| 9130-9134  | [konnectd](https://github.com/owncloud/ocis/tree/master/konnectd)             |
| 9135-9139  | [graph-explorer](https://github.com/owncloud/ocis/tree/master/graph-explorer) |
| 9140-9179  | [reva/storage](https://github.com/owncloud/ocis/tree/master/storage)          |
| 9180-9184  | FREE (formerly used by accounts)                                              |
| 9185-9189  | [thumbnails](https://github.com/owncloud/ocis/tree/master/thumbnails)         |
| 9190-9194  | [settings](https://github.com/owncloud/ocis/tree/master/settings)             |
| 9195-9199  | [store](https://github.com/owncloud/ocis/tree/master/store)                   |
| 9200-9204  | [proxy](https://github.com/owncloud/ocis/tree/master/proxy)                   |
| 9205-9209  | [markdown-editor](https://github.com/owncloud/ocis-markdown-editor)           |
| 9210-9214  | [reva](https://github.com/owncloud/ocis-reva) unused?                         |
| 9215-9219  | reva metadata storage                                                         |
| 9220-9224  | search                                                                          |
| 9225-9229  | photoprism (state: PoC)                                                       |
| 9230-9234  | [nats](https://github.com/owncloud/ocis/tree/master/nats)                     |
| 9235-9239  | [idm](https://github.com/owncloud/ocis/tree/master/proxy)                     |
| 9240-9244  | [app-registry](https://github.com/owncloud/ocis/tree/master/extensions/app-registry)     |
| 9245-9249  | FREE                                                                          |
| 9250-9254  | oCIS Runtime                                                                  |
| 9255-9259  | FREE                                                                          |
| 9260-9264  | FREE                                                                          |
| 9265-9269  | FREE                                                                          |
| 9270-9274  | FREE                                                                          |
| 9275-9279  | FREE                                                                          |
| 9280-9284  | FREE                                                                          |
| 9285-9289  | FREE                                                                          |
| 9290-9294  | FREE                                                                          |
| 9295-9299  | FREE                                                                          |
| 9300-9304  | FREE                                                                          |
| 9305-9309  | FREE                                                                          |
| 9310-9314  | FREE                                                                          |
| 9315-9319  | FREE                                                                          |
| 9320-9324  | FREE                                                                          |
| 9325-9329  | FREE                                                                          |
| 9330-9334  | FREE                                                                          |
| 9335-9339  | FREE                                                                          |
| 9340-9344  | FREE                                                                          |
| 9345-9349  | FREE                                                                          |
| 9350-9354  | FREE                                                                          |
| 9355-9359  | FREE                                                                          |
| 9360-9364  | FREE                                                                          |
| 9365-9369  | FREE                                                                          |
| 9370-9374  | FREE                                                                          |
| 9375-9379  | FREE                                                                          |
| 9380-9384  | FREE                                                                          |
| 9385-9389  | FREE                                                                          |
| 9390-9394  | FREE                                                                          |
| 9395-9399  | FREE                                                                          |
| 9400-9404  | FREE                                                                          |
| 9405-9409  | FREE                                                                          |
| 9410-9414  | FREE                                                                          |
| 9415-9419  | FREE                                                                          |
| 9420-9424  | FREE                                                                          |
| 9425-9429  | FREE                                                                          |
| 9430-9434  | FREE                                                                          |
| 9435-9439  | FREE                                                                          |
| 9440-9444  | FREE                                                                          |
| 9445-9449  | FREE                                                                          |
| 9450-9454  | FREE                                                                          |
| 9455-9459  | FREE                                                                          |
| 9460-9464  | FREE                                                                          |
| 9465-9469  | FREE                                                                          |
| 9470-9474  | FREE                                                                          |
| 9475-9479  | FREE                                                                          |
| 9480-9484  | FREE                                                                          |
| 9485-9489  | FREE                                                                          |
| 9490-9494  | FREE                                                                          |
| 9495-9499  | FREE                                                                          |
| 9500-9504  | FREE                                                                          |
| 9505-9509  | FREE                                                                          |
| 9510-9514  | FREE                                                                          |
| 9515-9519  | FREE                                                                          |
| 9520-9524  | FREE                                                                          |
| 9525-9529  | FREE                                                                          |
| 9530-9534  | FREE                                                                          |
| 9535-9539  | FREE                                                                          |
| 9540-9544  | FREE                                                                          |
| 9545-9549  | FREE                                                                          |
| 9550-9554  | FREE                                                                          |
| 9555-9559  | FREE                                                                          |
| 9560-9564  | FREE                                                                          |
| 9565-9569  | FREE                                                                          |
| 9570-9574  | FREE                                                                          |
| 9575-9579  | FREE                                                                          |
| 9580-9584  | FREE                                                                          |
| 9585-9589  | FREE                                                                          |
| 9590-9594  | FREE                                                                          |
| 9595-9599  | FREE                                                                          |
| 9600-9604  | FREE                                                                          |
| 9605-9609  | FREE                                                                          |
| 9610-9614  | FREE                                                                          |
| 9615-9619  | FREE                                                                          |
| 9620-9624  | FREE                                                                          |
| 9625-9629  | FREE                                                                          |
| 9630-9634  | FREE                                                                          |
| 9635-9639  | FREE                                                                          |
| 9640-9644  | FREE                                                                          |
| 9645-9649  | FREE                                                                          |
| 9650-9654  | FREE                                                                          |
| 9655-9659  | FREE                                                                          |
| 9660-9664  | FREE                                                                          |
| 9665-9669  | FREE                                                                          |
| 9670-9674  | FREE                                                                          |
| 9675-9679  | FREE                                                                          |
| 9680-9684  | FREE                                                                          |
| 9685-9689  | FREE                                                                          |
| 9690-9694  | FREE                                                                          |
| 9695-9699  | FREE                                                                          |
| 9700-9704  | FREE                                                                          |
| 9705-9709  | FREE                                                                          |
| 9710-9714  | FREE                                                                          |
| 9715-9719  | FREE                                                                          |
| 9720-9724  | FREE                                                                          |
| 9725-9729  | FREE                                                                          |
| 9730-9734  | FREE                                                                          |
| 9735-9739  | FREE                                                                          |
| 9740-9744  | FREE                                                                          |
| 9745-9749  | FREE                                                                          |
| 9750-9754  | FREE                                                                          |
| 9755-9759  | FREE                                                                          |
| 9760-9764  | FREE                                                                          |
| 9765-9769  | FREE                                                                          |
| 9770-9774  | FREE                                                                          |
| 9775-9779  | FREE                                                                          |
| 9780-9784  | FREE                                                                          |
| 9785-9789  | FREE                                                                          |
| 9790-9794  | FREE                                                                          |
| 9795-9799  | FREE                                                                          |
| 9800-9804  | FREE                                                                          |
| 9805-9809  | FREE                                                                          |
| 9810-9814  | FREE                                                                          |
| 9815-9819  | FREE                                                                          |
| 9820-9824  | FREE                                                                          |
| 9825-9829  | FREE                                                                          |
| 9830-9834  | FREE                                                                          |
| 9835-9839  | FREE                                                                          |
| 9840-9844  | FREE                                                                          |
| 9845-9849  | FREE                                                                          |
| 9850-9854  | FREE                                                                          |
| 9855-9859  | FREE                                                                          |
| 9860-9864  | FREE                                                                          |
| 9865-9869  | FREE                                                                          |
| 9870-9874  | FREE                                                                          |
| 9875-9879  | FREE                                                                          |
| 9880-9884  | FREE                                                                          |
| 9885-9889  | FREE                                                                          |
| 9890-9894  | FREE                                                                          |
| 9895-9899  | FREE                                                                          |
| 9900-9904  | FREE                                                                          |
| 9905-9909  | FREE                                                                          |
| 9910-9914  | FREE                                                                          |
| 9915-9919  | FREE                                                                          |
| 9920-9924  | FREE                                                                          |
| 9925-9929  | FREE                                                                          |
| 9930-9934  | FREE                                                                          |
| 9935-9939  | FREE                                                                          |
| 9940-9944  | FREE                                                                          |
| 9945-9949  | FREE                                                                          |
| 9950-9954  | FREE                                                                          |
| 9955-9959  | FREE                                                                          |
| 9960-9964  | FREE                                                                          |
| 9965-9969  | FREE                                                                          |
| 9970-9974  | FREE                                                                          |
| 9975-9979  | FREE                                                                          |
| 9980-9984  | FREE                                                                          |
| 9985-9989  | FREE                                                                          |
| 9990-9994  | FREE                                                                          |
| 9995-9999  | FREE                                                                          |
