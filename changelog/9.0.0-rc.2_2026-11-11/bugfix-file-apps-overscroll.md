Bugfix: Disable rubber-band overscroll of the page in file apps

Scrolling to the top or bottom of the content in a file app such as
OnlyOffice, Collabora or the text editor could make the whole page bounce
along with it, including the top bar and the app frame, briefly revealing the
blank page background behind them. The vertical elastic overscroll of the page
is now disabled while a file is open in one of these apps, so the surrounding
layout stays in place.

https://github.com/owncloud/ocis/pull/13140
