Enhancement: Remove unnecessary Service.Init()

As it turns out oCIS already calls this method. Invoking it twice would end in accidentally resetting values.

https://github.com/owncloud/ocis/pull/1705
