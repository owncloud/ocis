/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package provider

import (
	"html/template"
)

var checkSessionIframeTemplate = template.Must(template.New("check-session.html").Parse(`
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<script id="cookie-name" type="application/x-kkbs-name">{{.CookieName}}</script>
<script type="text/javascript" nonce={{.Nonce}}>
/*
Forge-SHA256 - https://github.com/brillout/forge-sha256/tree/6ad5535e0be2385fdc53f1d9ce2b172365c70333

The MIT License (MIT)

Copyright (c) 2015-2017 Romuald Brillout and contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Forge project - https://github.com/digitalbazaar/forge

New BSD License (3-clause)
Copyright (c) 2010, Digital Bazaar, Inc.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of Digital Bazaar, Inc. nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL DIGITAL BAZAAR BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
(function(){function p(a){this.data="";this.a=0;if("string"===typeof a)this.data=a;else if(b.D(a)||b.L(a)){a=new Uint8Array(a);try{this.data=String.fromCharCode.apply(null,a)}catch(f){for(var v=0;v<a.length;++v)this.M(a[v])}}else if(a instanceof p||"object"===typeof a&&"string"===typeof a.data&&"number"===typeof a.a)this.data=a.data,this.a=a.a;this.v=0}function w(a,f,b){for(var d,c,h,m,g,k,e,r,n,l,t,q,u,p=b.length();64<=p;){for(g=0;16>g;++g)f[g]=b.getInt32();for(;64>g;++g)d=f[g-2],d=(d>>>17|d<<15)^
(d>>>19|d<<13)^d>>>10,c=f[g-15],c=(c>>>7|c<<25)^(c>>>18|c<<14)^c>>>3,f[g]=d+f[g-7]+c+f[g-16]|0;k=a.g;e=a.h;r=a.i;n=a.j;l=a.l;t=a.m;q=a.o;u=a.s;for(g=0;64>g;++g)d=(l>>>6|l<<26)^(l>>>11|l<<21)^(l>>>25|l<<7),h=q^l&(t^q),c=(k>>>2|k<<30)^(k>>>13|k<<19)^(k>>>22|k<<10),m=k&e|r&(k^e),d=u+d+h+x[g]+f[g],c+=m,u=q,q=t,t=l,l=n+d|0,n=r,r=e,e=k,k=d+c|0;a.g=a.g+k|0;a.h=a.h+e|0;a.i=a.i+r|0;a.j=a.j+n|0;a.l=a.l+l|0;a.m=a.m+t|0;a.o=a.o+q|0;a.s=a.s+u|0;p-=64}}var m,y,e,b=m=m||{};b.D=function(a){return"undefined"!==typeof ArrayBuffer&&
a instanceof ArrayBuffer};b.L=function(a){return a&&b.D(a.buffer)&&void 0!==a.byteLength};b.G=p;b.b=p;b.b.prototype.H=function(a){this.v+=a;4096<this.v&&(this.v=0)};b.b.prototype.length=function(){return this.data.length-this.a};b.b.prototype.M=function(a){this.u(String.fromCharCode(a))};b.b.prototype.u=function(a){this.data+=a;this.H(a.length)};b.b.prototype.c=function(a){this.u(String.fromCharCode(a>>24&255)+String.fromCharCode(a>>16&255)+String.fromCharCode(a>>8&255)+String.fromCharCode(a&255))};
b.b.prototype.getInt16=function(){var a=this.data.charCodeAt(this.a)<<8^this.data.charCodeAt(this.a+1);this.a+=2;return a};b.b.prototype.getInt32=function(){var a=this.data.charCodeAt(this.a)<<24^this.data.charCodeAt(this.a+1)<<16^this.data.charCodeAt(this.a+2)<<8^this.data.charCodeAt(this.a+3);this.a+=4;return a};b.b.prototype.B=function(){return this.data.slice(this.a)};b.b.prototype.compact=function(){0<this.a&&(this.data=this.data.slice(this.a),this.a=0);return this};b.b.prototype.clear=function(){this.data=
"";this.a=0;return this};b.b.prototype.truncate=function(a){a=Math.max(0,this.length()-a);this.data=this.data.substr(this.a,a);this.a=0;return this};b.b.prototype.N=function(){for(var a="",f=this.a;f<this.data.length;++f){var b=this.data.charCodeAt(f);16>b&&(a+="0");a+=b.toString(16)}return a};b.b.prototype.toString=function(){return b.I(this.B())};b.createBuffer=function(a,f){void 0!==a&&"utf8"===(f||"raw")&&(a=b.C(a));return new b.G(a)};b.J=function(){for(var a=String.fromCharCode(0),b=64,e="";0<
b;)b&1&&(e+=a),b>>>=1,0<b&&(a+=a);return e};b.C=function(a){return unescape(encodeURIComponent(a))};b.I=function(a){return decodeURIComponent(escape(a))};b.K=function(a){for(var b=0;b<a.length;b++)if(a.charCodeAt(b)>>>8)return!0;return!1};var z=y=y||{};e=e||{};e.A=e.A||{};e.F=e.A.F=z;z.create=function(){A||(n=String.fromCharCode(128),n+=m.J(),x=[1116352408,1899447441,3049323471,3921009573,961987163,1508970993,2453635748,2870763221,3624381080,310598401,607225278,1426881987,1925078388,2162078206,2614888103,
3248222580,3835390401,4022224774,264347078,604807628,770255983,1249150122,1555081692,1996064986,2554220882,2821834349,2952996808,3210313671,3336571891,3584528711,113926993,338241895,666307205,773529912,1294757372,1396182291,1695183700,1986661051,2177026350,2456956037,2730485921,2820302411,3259730800,3345764771,3516065817,3600352804,4094571909,275423344,430227734,506948616,659060556,883997877,958139571,1322822218,1537002063,1747873779,1955562222,2024104815,2227730452,2361852424,2428436474,2756734187,
3204031479,3329325298],A=!0);var a=null,b=m.createBuffer(),e=Array(64),d={algorithm:"sha256",O:64,P:32,w:0,f:[0,0],start:function(){d.w=0;d.f=[0,0];b=m.createBuffer();a={g:1779033703,h:3144134277,i:1013904242,j:2773480762,l:1359893119,m:2600822924,o:528734635,s:1541459225};return d}};d.start();d.update=function(c,h){"utf8"===h&&(c=m.C(c));d.w+=c.length;d.f[0]+=c.length/4294967296>>>0;d.f[1]+=c.length>>>0;b.u(c);w(a,e,b);(2048<b.a||0===b.length())&&b.compact();return d};d.digest=function(){var c=m.createBuffer();
c.u(b.B());c.u(n.substr(0,64-(d.f[1]+8&63)));c.c(d.f[0]<<3|d.f[0]>>>28);c.c(d.f[1]<<3);var h={g:a.g,h:a.h,i:a.i,j:a.j,l:a.l,m:a.m,o:a.o,s:a.s};w(h,e,c);c=m.createBuffer();c.c(h.g);c.c(h.h);c.c(h.i);c.c(h.j);c.c(h.l);c.c(h.m);c.c(h.o);c.c(h.s);return c};return d};var n=null,A=!1,x=null;window.forge_sha256=function(a){var f=e.F.create();f.update(a,b.K(a)?"utf8":void 0);return f.digest().N()}})();
</script>
</head>
<body>
<noscript>
  You need to enable JavaScript to run this app.
</noscript>
<script type="text/javascript" nonce={{.Nonce}}>
// This implements OpenID Connect Session Managament 1.0 OP IFrame as specified
// in https://openid.net/specs/openid-connect-session-1_0.html#OPiframe
(function() {
	var cache = {};

	// Get cookie name.
	var cookieNameElem = document.getElementById('cookie-name');
	if (cookieNameElem) {
		var cookieName = cookieNameElem.textContent.trim();
	}

	// Helper to get cookie by name.
	function getCookie(name) {
		var value = '; ' + document.cookie;
		var parts = value.split('; ' + name + "=");
		if (parts.length === 2) {
			return parts.pop().split(';').shift();
		}
	}

	function getBrowserState() {
		var browserState = getCookie(cookieName);
		return browserState ? browserState : '';
	}

	function makeSessionState(clientID, origin, browserState, salt) {
		return ''+forge_sha256(clientID + ' ' + origin + ' ' + browserState + ' ' + salt);
	}

	// Outer error catching data processor.
	function receiveMessage(origin, data) {
		if (!origin || !data) {
			return 'error';
		}

		try {
			return handleMessage(origin, data);
		} catch (err) {
			return 'error';
		}
	}

	// Parse and validatie.
	function handleMessage(origin, data) {
		// Parse data.
		var dataParts = data.split(' ');
		if (dataParts.length !== 2) {
			return 'error';
		}
		var clientID = dataParts[0];
		var sessionState = dataParts[1];
		if (!clientID || !sessionState) {
			return 'error';
		}
		var sessionStateParts = sessionState.split('.');
		if (sessionStateParts.length != 2) {
			return 'error';
		}
		var clientStateHash = sessionStateParts[0];
		var salt = sessionStateParts[1];
		if (!clientStateHash || !salt) {
			return 'error';
		}

		// Get browser state.
		var browserState = getBrowserState();

		// Make session state.
		var expectedStateHash = makeSessionState(clientID, origin, browserState, salt);

		// Compare.
		var status = clientStateHash === expectedStateHash ? 'unchanged' : 'changed';

		// Cache.
		if (cache.browserState === browserState && status === 'changed' && cache.status === status) {
			// If the browser state remains unchanged, something with the cookie
			// did not work properly. To avoid a fast repeating loop, it is
			// better to fail with error.
			return 'error';
		}
		cache.status = status;
		cache.browserState = browserState;

		return status;
	}

	// Register event.
	if (cookieName && window.parent !== window) {
		window.addEventListener('message', function(event) {
			// Only do something when receiving a message from our parent or
			// from another window which shares our parent.
			if (window.parent === event.source || window.parent === event.source.parent) {
				var response = receiveMessage(event.origin, event.data);
				event.source.postMessage(response, event.origin);
			}
		}, false);
	}
})();
</script>
</body>
</html>
`))
