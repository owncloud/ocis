Enhancement: allow querying federated user roles for sharing

When listing permissions clients can now fetch the list of available federated sharing roles by sending a `GET /graph/v1beta1/drives/{driveid}/items/{itemid}/permissions?$filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition, '@Subject.UserType=="Federated"'))` request. Note that this is the only supported filter expression. Federated sharing roles will be omitted from requests without this filter.

https://github.com/owncloud/ocis/pull/9765
