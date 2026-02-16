Enhancement: Added generic way to translate composite entities

Added a generic way to translate the necessary fields in composite entities.
The function takes the entity, translation function and fields to translate that are described by the TranslateField function.
The function supports nested structs and slices of structs.

https://github.com/owncloud/ocis/pull/9722
https://github.com/owncloud/ocis/issues/9700
