Bugfix: Roll back tag assignment when the event cannot be published

Assigning a tag via `PUT /graph/v1.0/extensions/org.libregraph/tags` answered
200 even when the `TagsAdded` event could not be published. The file carried the
tag, but the search index never learned about it, so the tag was missing from the
tag list and from `tag:` searches.

The graph service now restores the previous tags and answers 500 in that case,
as it already does when removing tags.

https://github.com/owncloud/ocis/issues/13286
