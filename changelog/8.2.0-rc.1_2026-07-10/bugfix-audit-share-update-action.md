Bugfix: Always set an audit action for share updates

Share-update audit events were written with an empty `action` and a message of
`updated field ''`, because the conversion read the deprecated
`ShareUpdated.Updated` field, which reva no longer populates (it now sets
`UpdateMask`). Received-share declines were also not audited correctly because
the conversion matched `SHARE_STATE_DECLINED`, which is not a CS3 share state
(the enum value is `SHARE_STATE_REJECTED`), again producing an empty action and
message.

The conversion now derives the updated field and action from `UpdateMask`
(falling back to the deprecated field), maps `SHARE_STATE_REJECTED` to the
declined action, and uses a generic, non-empty action for any unrecognized
update field or share state, so every audit entry carries a meaningful,
countable action.

https://github.com/owncloud/ocis/issues/7661
https://github.com/owncloud/ocis/pull/12423
