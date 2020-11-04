Bugfix: Allow consent-prompt with switch-account

Multiple prompt values are allowed and this change fixes the check for
select_account if it was used together with other prompt values. Where
select_account previously was ignored, it is now processed as required,
fixing the use case when a RP wants to trigger select_account first
while at the same time wants also to request interactive consent.

https://github.com/owncloud/ocis/pull/788
