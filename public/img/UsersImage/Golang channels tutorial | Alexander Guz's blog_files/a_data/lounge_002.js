!function(){"use strict";var a=window.document,b={STYLES:"https://c.disquscdn.com/next/embed/styles/lounge.91c71242b4acaa0ee7f9db125ef21f90.css",RTL_STYLES:"https://c.disquscdn.com/next/embed/styles/lounge_rtl.2cdf4166ca99f0eb2a97e8de75583afb.css","lounge/main":"https://c.disquscdn.com/next/embed/lounge.bundle.0523614fef9787c3e3459e0602078385.js","discovery/main":"https://c.disquscdn.com/next/embed/discovery.bundle.a6a4c220420629ec61e18a6cd65bf4e1.js","recommendations/main":"https://c.disquscdn.com/next/embed/recommendations.bundle.9b1ee01767aef897cc2467236bb0313b.js","remote/config":"https://disqus.com/next/config.js","common/vendor_extensions/highlight":"https://c.disquscdn.com/next/embed/highlight.6fbf348532f299e045c254c49c4dbedf.js"};window.require={baseUrl:"https://c.disquscdn.com/next/current/embed",paths:["lounge/main","discovery/main","recommendations/main","remote/config","common/vendor_extensions/highlight"].reduce(function(a,c){return a[c]=b[c].slice(0,-3),a},{})};var c=a.createElement("script");c.onload=function(){require(["common/main"],function(a){a.init("lounge",b)})},c.src="https://c.disquscdn.com/next/embed/common.bundle.0b9bbdb3bc568241a5d1d7626947e8b0.js",a.body.appendChild(c)}();