var qe=Object.defineProperty;var Ke=(e,t,n)=>t in e?qe(e,t,{enumerable:!0,configurable:!0,writable:!0,value:n}):e[t]=n;var X=(e,t,n)=>(Ke(e,typeof t!="symbol"?t+"":t,n),n);(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const r of document.querySelectorAll('link[rel="modulepreload"]'))s(r);new MutationObserver(r=>{for(const i of r)if(i.type==="childList")for(const o of i.addedNodes)o.tagName==="LINK"&&o.rel==="modulepreload"&&s(o)}).observe(document,{childList:!0,subtree:!0});function n(r){const i={};return r.integrity&&(i.integrity=r.integrity),r.referrerPolicy&&(i.referrerPolicy=r.referrerPolicy),r.crossOrigin==="use-credentials"?i.credentials="include":r.crossOrigin==="anonymous"?i.credentials="omit":i.credentials="same-origin",i}function s(r){if(r.ep)return;r.ep=!0;const i=n(r);fetch(r.href,i)}})();const m={context:void 0,registry:void 0};function H(e){m.context=e}const Ge=(e,t)=>e===t,Xe=Symbol("solid-track"),Q={equals:Ge};let ve=Me;const z=1,Y=2,ke={owned:null,cleanups:null,context:null,owner:null},le={};var w=null;let oe=null,p=null,_=null,M=null,ee=0;function F(e,t){const n=p,s=w,r=e.length===0,i=t===void 0?s:t,o=r?ke:{owned:null,cleanups:null,context:i?i.context:null,owner:i},l=r?e:()=>e(()=>v(()=>ne(o)));w=o,p=null;try{return R(l,!0)}finally{p=n,w=s}}function O(e,t){t=t?Object.assign({},Q,t):Q;const n={value:e,observers:null,observerSlots:null,comparator:t.equals||void 0},s=r=>(typeof r=="function"&&(r=r(n.value)),Ie(n,r));return[Oe.bind(n),s]}function we(e,t,n){const s=te(e,t,!0,z);B(s)}function j(e,t,n){const s=te(e,t,!1,z);B(s)}function he(e,t,n){ve=nt;const s=te(e,t,!1,z);(!n||!n.render)&&(s.user=!0),M?M.push(s):B(s)}function I(e,t,n){n=n?Object.assign({},Q,n):Q;const s=te(e,t,!0,0);return s.observers=null,s.observerSlots=null,s.comparator=n.equals||void 0,B(s),Oe.bind(s)}function Qe(e){return e&&typeof e=="object"&&"then"in e}function Ye(e,t,n){let s,r,i;arguments.length===2&&typeof t=="object"||arguments.length===1?(s=!0,r=e,i=t||{}):(s=e,r=t,i=n||{});let o=null,l=le,u=null,f=!1,c="initialValue"in i,d=typeof s=="function"&&I(s);const a=new Set,[y,E]=(i.storage||O)(i.initialValue),[g,N]=O(void 0),[k,C]=O(void 0,{equals:!1}),[L,T]=O(c?"ready":"unresolved");if(m.context){u=`${m.context.id}${m.context.count++}`;let x;i.ssrLoadFrom==="initial"?l=i.initialValue:m.load&&(x=m.load(u))&&(l=x)}function P(x,$,S,D){return o===x&&(o=null,D!==void 0&&(c=!0),(x===l||$===l)&&i.onHydrated&&queueMicrotask(()=>i.onHydrated(D,{value:$})),l=le,se($,S)),$}function se(x,$){R(()=>{$===void 0&&E(()=>x),T($!==void 0?"errored":c?"ready":"unresolved"),N($);for(const S of a.keys())S.decrement();a.clear()},!1)}function re(){const x=Ze,$=y(),S=g();if(S!==void 0&&!o)throw S;return p&&!p.user&&x&&we(()=>{k(),o&&(x.resolved||a.has(x)||(x.increment(),a.add(x)))}),$}function ie(x=!0){if(x!==!1&&f)return;f=!1;const $=d?d():s;if($==null||$===!1){P(o,v(y));return}const S=l!==le?l:v(()=>r($,{value:y(),refetching:x}));return Qe(S)?(o=S,"value"in S?(S.status==="success"?P(o,S.value,void 0,$):P(o,void 0,void 0,$),S):(f=!0,queueMicrotask(()=>f=!1),R(()=>{T(c?"refreshing":"pending"),C()},!1),S.then(D=>P(S,D,void 0,$),D=>P(S,void 0,ze(D),$)))):(P(o,S,void 0,$),S)}return Object.defineProperties(re,{state:{get:()=>L()},error:{get:()=>g()},loading:{get(){const x=L();return x==="pending"||x==="refreshing"}},latest:{get(){if(!c)return re();const x=g();if(x&&!o)throw x;return y()}}}),d?we(()=>ie(!1)):ie(!1),[re,{refetch:ie,mutate:E}]}function v(e){if(p===null)return e();const t=p;p=null;try{return e()}finally{p=t}}function on(e,t,n){const s=Array.isArray(e);let r,i=n&&n.defer;return o=>{let l;if(s){l=Array(e.length);for(let f=0;f<e.length;f++)l[f]=e[f]()}else l=e();if(i){i=!1;return}const u=v(()=>t(l,r,o));return r=l,u}}function Ne(e){he(()=>v(e))}function q(e){return w===null||(w.cleanups===null?w.cleanups=[e]:w.cleanups.push(e)),e}function ye(){return w}function Je(e,t){const n=w,s=p;w=e,p=null;try{return R(t,!0)}catch(r){me(r)}finally{w=n,p=s}}function Pe(e,t){const n=Symbol("context");return{id:n,Provider:st(n),defaultValue:e}}function pe(e){return w&&w.context&&w.context[e.id]!==void 0?w.context[e.id]:e.defaultValue}function Te(e){const t=I(e),n=I(()=>ae(t()));return n.toArray=()=>{const s=n();return Array.isArray(s)?s:s!=null?[s]:[]},n}let Ze;function Oe(){if(this.sources&&this.state)if(this.state===z)B(this);else{const e=_;_=null,R(()=>Z(this),!1),_=e}if(p){const e=this.observers?this.observers.length:0;p.sources?(p.sources.push(this),p.sourceSlots.push(e)):(p.sources=[this],p.sourceSlots=[e]),this.observers?(this.observers.push(p),this.observerSlots.push(p.sources.length-1)):(this.observers=[p],this.observerSlots=[p.sources.length-1])}return this.value}function Ie(e,t,n){let s=e.value;return(!e.comparator||!e.comparator(s,t))&&(e.value=t,e.observers&&e.observers.length&&R(()=>{for(let r=0;r<e.observers.length;r+=1){const i=e.observers[r],o=oe&&oe.running;o&&oe.disposed.has(i),(o?!i.tState:!i.state)&&(i.pure?_.push(i):M.push(i),i.observers&&je(i)),o||(i.state=z)}if(_.length>1e6)throw _=[],new Error},!1)),t}function B(e){if(!e.fn)return;ne(e);const t=ee;et(e,e.value,t)}function et(e,t,n){let s;const r=w,i=p;p=w=e;try{s=e.fn(t)}catch(o){return e.pure&&(e.state=z,e.owned&&e.owned.forEach(ne),e.owned=null),e.updatedAt=n+1,me(o)}finally{p=i,w=r}(!e.updatedAt||e.updatedAt<=n)&&(e.updatedAt!=null&&"observers"in e?Ie(e,s):e.value=s,e.updatedAt=n)}function te(e,t,n,s=z,r){const i={fn:e,state:s,updatedAt:null,owned:null,sources:null,sourceSlots:null,cleanups:null,value:t,owner:w,context:w?w.context:null,pure:n};return w===null||w!==ke&&(w.owned?w.owned.push(i):w.owned=[i]),i}function J(e){if(e.state===0)return;if(e.state===Y)return Z(e);if(e.suspense&&v(e.suspense.inFallback))return e.suspense.effects.push(e);const t=[e];for(;(e=e.owner)&&(!e.updatedAt||e.updatedAt<ee);)e.state&&t.push(e);for(let n=t.length-1;n>=0;n--)if(e=t[n],e.state===z)B(e);else if(e.state===Y){const s=_;_=null,R(()=>Z(e,t[0]),!1),_=s}}function R(e,t){if(_)return e();let n=!1;t||(_=[]),M?n=!0:M=[],ee++;try{const s=e();return tt(n),s}catch(s){n||(M=null),_=null,me(s)}}function tt(e){if(_&&(Me(_),_=null),e)return;const t=M;M=null,t.length&&R(()=>ve(t),!1)}function Me(e){for(let t=0;t<e.length;t++)J(e[t])}function nt(e){let t,n=0;for(t=0;t<e.length;t++){const s=e[t];s.user?e[n++]=s:J(s)}if(m.context){if(m.count){m.effects||(m.effects=[]),m.effects.push(...e.slice(0,n));return}else m.effects&&(e=[...m.effects,...e],n+=m.effects.length,delete m.effects);H()}for(t=0;t<n;t++)J(e[t])}function Z(e,t){e.state=0;for(let n=0;n<e.sources.length;n+=1){const s=e.sources[n];if(s.sources){const r=s.state;r===z?s!==t&&(!s.updatedAt||s.updatedAt<ee)&&J(s):r===Y&&Z(s,t)}}}function je(e){for(let t=0;t<e.observers.length;t+=1){const n=e.observers[t];n.state||(n.state=Y,n.pure?_.push(n):M.push(n),n.observers&&je(n))}}function ne(e){let t;if(e.sources)for(;e.sources.length;){const n=e.sources.pop(),s=e.sourceSlots.pop(),r=n.observers;if(r&&r.length){const i=r.pop(),o=n.observerSlots.pop();s<r.length&&(i.sourceSlots[o]=s,r[s]=i,n.observerSlots[s]=o)}}if(e.owned){for(t=e.owned.length-1;t>=0;t--)ne(e.owned[t]);e.owned=null}if(e.cleanups){for(t=e.cleanups.length-1;t>=0;t--)e.cleanups[t]();e.cleanups=null}e.state=0}function ze(e){return e instanceof Error?e:new Error(typeof e=="string"?e:"Unknown error",{cause:e})}function me(e,t=w){throw ze(e)}function ae(e){if(typeof e=="function"&&!e.length)return ae(e());if(Array.isArray(e)){const t=[];for(let n=0;n<e.length;n++){const s=ae(e[n]);Array.isArray(s)?t.push.apply(t,s):t.push(s)}return t}return e}function st(e,t){return function(s){let r;return j(()=>r=v(()=>(w.context={...w.context,[e]:s.value},Te(()=>s.children))),void 0),r}}const rt=Symbol("fallback");function be(e){for(let t=0;t<e.length;t++)e[t]()}function it(e,t,n={}){let s=[],r=[],i=[],o=0,l=t.length>1?[]:null;return q(()=>be(i)),()=>{let u=e()||[],f,c;return u[Xe],v(()=>{let a=u.length,y,E,g,N,k,C,L,T,P;if(a===0)o!==0&&(be(i),i=[],s=[],r=[],o=0,l&&(l=[])),n.fallback&&(s=[rt],r[0]=F(se=>(i[0]=se,n.fallback())),o=1);else if(o===0){for(r=new Array(a),c=0;c<a;c++)s[c]=u[c],r[c]=F(d);o=a}else{for(g=new Array(a),N=new Array(a),l&&(k=new Array(a)),C=0,L=Math.min(o,a);C<L&&s[C]===u[C];C++);for(L=o-1,T=a-1;L>=C&&T>=C&&s[L]===u[T];L--,T--)g[T]=r[L],N[T]=i[L],l&&(k[T]=l[L]);for(y=new Map,E=new Array(T+1),c=T;c>=C;c--)P=u[c],f=y.get(P),E[c]=f===void 0?-1:f,y.set(P,c);for(f=C;f<=L;f++)P=s[f],c=y.get(P),c!==void 0&&c!==-1?(g[c]=r[f],N[c]=i[f],l&&(k[c]=l[f]),c=E[c],y.set(P,c)):i[f]();for(c=C;c<a;c++)c in g?(r[c]=g[c],i[c]=N[c],l&&(l[c]=k[c],l[c](c))):r[c]=F(d);r=r.slice(0,o=a),s=u.slice(0)}return r});function d(a){if(i[c]=a,l){const[y,E]=O(c);return l[c]=E,t(u[c],y)}return t(u[c])}}}function h(e,t){return v(()=>e(t||{}))}function lt(e){let t,n;const s=r=>{const i=m.context;if(i){const[l,u]=O();m.count||(m.count=0),m.count++,(n||(n=e())).then(f=>{H(i),m.count--,u(()=>f.default),H()}),t=l}else if(!t){const[l]=Ye(()=>(n||(n=e())).then(u=>u.default));t=l}let o;return I(()=>(o=t())&&v(()=>{if(!i)return o(r);const l=m.context;H(i);const u=o(r);return H(l),u}))};return s.preload=()=>n||((n=e()).then(r=>t=()=>r.default),n),s}const Re=e=>`Stale read from <${e}>.`;function ot(e){const t="fallback"in e&&{fallback:()=>e.fallback};return I(it(()=>e.each,e.children,t||void 0))}function G(e){const t=e.keyed,n=I(()=>e.when,void 0,{equals:(s,r)=>t?s===r:!s==!r});return I(()=>{const s=n();if(s){const r=e.children;return typeof r=="function"&&r.length>0?v(()=>r(t?s:()=>{if(!v(n))throw Re("Show");return e.when})):r}return e.fallback},void 0,void 0)}function De(e){let t=!1;const n=(i,o)=>i[0]===o[0]&&(t?i[1]===o[1]:!i[1]==!o[1])&&i[2]===o[2],s=Te(()=>e.children),r=I(()=>{let i=s();Array.isArray(i)||(i=[i]);for(let o=0;o<i.length;o++){const l=i[o].when;if(l)return t=!!i[o].keyed,[o,l,i[o]]}return[-1]},void 0,{equals:n});return I(()=>{const[i,o,l]=r();if(i<0)return e.fallback;const u=l.children;return typeof u=="function"&&u.length>0?v(()=>u(t?o:()=>{if(v(r)[0]!==i)throw Re("Match");return l.when})):u},void 0,void 0)}function W(e){return e}const ct=["allowfullscreen","async","autofocus","autoplay","checked","controls","default","disabled","formnovalidate","hidden","indeterminate","inert","ismap","loop","multiple","muted","nomodule","novalidate","open","playsinline","readonly","required","reversed","seamless","selected"],ut=new Set(["className","value","readOnly","formNoValidate","isMap","noModule","playsInline",...ct]),ft=new Set(["innerHTML","textContent","innerText","children"]),at=Object.assign(Object.create(null),{className:"class",htmlFor:"for"}),dt=Object.assign(Object.create(null),{class:"className",formnovalidate:{$:"formNoValidate",BUTTON:1,INPUT:1},ismap:{$:"isMap",IMG:1},nomodule:{$:"noModule",SCRIPT:1},playsinline:{$:"playsInline",VIDEO:1},readonly:{$:"readOnly",INPUT:1,TEXTAREA:1}});function ht(e,t){const n=dt[e];return typeof n=="object"?n[t]?n.$:void 0:n}const mt=new Set(["beforeinput","click","dblclick","contextmenu","focusin","focusout","input","keydown","keyup","mousedown","mousemove","mouseout","mouseover","mouseup","pointerdown","pointermove","pointerout","pointerover","pointerup","touchend","touchmove","touchstart"]),gt={xlink:"http://www.w3.org/1999/xlink",xml:"http://www.w3.org/XML/1998/namespace"};function wt(e,t,n){let s=n.length,r=t.length,i=s,o=0,l=0,u=t[r-1].nextSibling,f=null;for(;o<r||l<i;){if(t[o]===n[l]){o++,l++;continue}for(;t[r-1]===n[i-1];)r--,i--;if(r===o){const c=i<s?l?n[l-1].nextSibling:n[i-l]:u;for(;l<i;)e.insertBefore(n[l++],c)}else if(i===l)for(;o<r;)(!f||!f.has(t[o]))&&t[o].remove(),o++;else if(t[o]===n[i-1]&&n[l]===t[r-1]){const c=t[--r].nextSibling;e.insertBefore(n[l++],t[o++].nextSibling),e.insertBefore(n[--i],c),t[r]=n[i]}else{if(!f){f=new Map;let d=l;for(;d<i;)f.set(n[d],d++)}const c=f.get(t[o]);if(c!=null)if(l<c&&c<i){let d=o,a=1,y;for(;++d<r&&d<i&&!((y=f.get(t[d]))==null||y!==c+a);)a++;if(a>c-l){const E=t[o];for(;l<c;)e.insertBefore(n[l++],E)}else e.replaceChild(n[l++],t[o++])}else o++;else t[o++].remove()}}}const xe="_$DX_DELEGATE";function yt(e,t,n,s={}){let r;return F(i=>{r=i,t===document?e():b(t,e(),t.firstChild?null:void 0,n)},s.owner),()=>{r(),t.textContent=""}}function A(e,t,n){let s;const r=()=>{const o=document.createElement("template");return o.innerHTML=e,n?o.content.firstChild.firstChild:o.content.firstChild},i=t?()=>v(()=>document.importNode(s||(s=r()),!0)):()=>(s||(s=r())).cloneNode(!0);return i.cloneNode=i,i}function ge(e,t=window.document){const n=t[xe]||(t[xe]=new Set);for(let s=0,r=e.length;s<r;s++){const i=e[s];n.has(i)||(n.add(i),t.addEventListener(i,Lt))}}function K(e,t,n){m.context||(n==null?e.removeAttribute(t):e.setAttribute(t,n))}function pt(e,t,n,s){m.context||(s==null?e.removeAttributeNS(t,n):e.setAttributeNS(t,n,s))}function bt(e,t){m.context||(t==null?e.removeAttribute("class"):e.className=t)}function Ue(e,t,n,s){if(s)Array.isArray(n)?(e[`$$${t}`]=n[0],e[`$$${t}Data`]=n[1]):e[`$$${t}`]=n;else if(Array.isArray(n)){const r=n[0];e.addEventListener(t,n[0]=i=>r.call(e,n[1],i))}else e.addEventListener(t,n)}function xt(e,t,n={}){const s=Object.keys(t||{}),r=Object.keys(n);let i,o;for(i=0,o=r.length;i<o;i++){const l=r[i];!l||l==="undefined"||t[l]||(Ee(e,l,!1),delete n[l])}for(i=0,o=s.length;i<o;i++){const l=s[i],u=!!t[l];!l||l==="undefined"||n[l]===u||!u||(Ee(e,l,!0),n[l]=u)}return n}function Et(e,t,n){if(!t)return n?K(e,"style"):t;const s=e.style;if(typeof t=="string")return s.cssText=t;typeof n=="string"&&(s.cssText=n=void 0),n||(n={}),t||(t={});let r,i;for(i in n)t[i]==null&&s.removeProperty(i),delete n[i];for(i in t)r=t[i],r!==n[i]&&(s.setProperty(i,r),n[i]=r);return n}function $t(e,t={},n,s){const r={};return s||j(()=>r.children=V(e,t.children,r.children)),j(()=>t.ref&&t.ref(e)),j(()=>St(e,t,n,!0,r,!0)),r}function At(e,t,n){return v(()=>e(t,n))}function b(e,t,n,s){if(n!==void 0&&!s&&(s=[]),typeof t!="function")return V(e,t,s,n);j(r=>V(e,t(),r,n),s)}function St(e,t,n,s,r={},i=!1){t||(t={});for(const o in r)if(!(o in t)){if(o==="children")continue;r[o]=$e(e,o,null,r[o],n,i)}for(const o in t){if(o==="children"){s||V(e,t.children);continue}const l=t[o];r[o]=$e(e,o,l,r[o],n,i)}}function Ct(e){return e.toLowerCase().replace(/-([a-z])/g,(t,n)=>n.toUpperCase())}function Ee(e,t,n){const s=t.trim().split(/\s+/);for(let r=0,i=s.length;r<i;r++)e.classList.toggle(s[r],n)}function $e(e,t,n,s,r,i){let o,l,u,f,c;if(t==="style")return Et(e,n,s);if(t==="classList")return xt(e,n,s);if(n===s)return s;if(t==="ref")i||n(e);else if(t.slice(0,3)==="on:"){const d=t.slice(3);s&&e.removeEventListener(d,s),n&&e.addEventListener(d,n)}else if(t.slice(0,10)==="oncapture:"){const d=t.slice(10);s&&e.removeEventListener(d,s,!0),n&&e.addEventListener(d,n,!0)}else if(t.slice(0,2)==="on"){const d=t.slice(2).toLowerCase(),a=mt.has(d);if(!a&&s){const y=Array.isArray(s)?s[0]:s;e.removeEventListener(d,y)}(a||n)&&(Ue(e,d,n,a),a&&ge([d]))}else if(t.slice(0,5)==="attr:")K(e,t.slice(5),n);else if((c=t.slice(0,5)==="prop:")||(u=ft.has(t))||!r&&((f=ht(t,e.tagName))||(l=ut.has(t)))||(o=e.nodeName.includes("-"))){if(c)t=t.slice(5),l=!0;else if(m.context)return n;t==="class"||t==="className"?bt(e,n):o&&!l&&!u?e[Ct(t)]=n:e[f||t]=n}else{const d=r&&t.indexOf(":")>-1&&gt[t.split(":")[0]];d?pt(e,d,t,n):K(e,at[t]||t,n)}return n}function Lt(e){const t=`$$${e.type}`;let n=e.composedPath&&e.composedPath()[0]||e.target;for(e.target!==n&&Object.defineProperty(e,"target",{configurable:!0,value:n}),Object.defineProperty(e,"currentTarget",{configurable:!0,get(){return n||document}}),m.registry&&!m.done&&(m.done=_$HY.done=!0);n;){const s=n[t];if(s&&!n.disabled){const r=n[`${t}Data`];if(r!==void 0?s.call(n,r,e):s.call(n,e),e.cancelBubble)return}n=n._$host||n.parentNode||n.host}}function V(e,t,n,s,r){if(m.context){!n&&(n=[...e.childNodes]);let l=[];for(let u=0;u<n.length;u++){const f=n[u];f.nodeType===8&&f.data.slice(0,2)==="!$"?f.remove():l.push(f)}n=l}for(;typeof n=="function";)n=n();if(t===n)return n;const i=typeof t,o=s!==void 0;if(e=o&&n[0]&&n[0].parentNode||e,i==="string"||i==="number"){if(m.context)return n;if(i==="number"&&(t=t.toString()),o){let l=n[0];l&&l.nodeType===3?l.data=t:l=document.createTextNode(t),n=U(e,n,s,l)}else n!==""&&typeof n=="string"?n=e.firstChild.data=t:n=e.textContent=t}else if(t==null||i==="boolean"){if(m.context)return n;n=U(e,n,s)}else{if(i==="function")return j(()=>{let l=t();for(;typeof l=="function";)l=l();n=V(e,l,n,s)}),()=>n;if(Array.isArray(t)){const l=[],u=n&&Array.isArray(n);if(de(l,t,n,r))return j(()=>n=V(e,l,n,s,!0)),()=>n;if(m.context){if(!l.length)return n;if(s===void 0)return[...e.childNodes];let f=l[0],c=[f];for(;(f=f.nextSibling)!==s;)c.push(f);return n=c}if(l.length===0){if(n=U(e,n,s),o)return n}else u?n.length===0?Ae(e,l,s):wt(e,n,l):(n&&U(e),Ae(e,l));n=l}else if(t.nodeType){if(m.context&&t.parentNode)return n=o?[t]:t;if(Array.isArray(n)){if(o)return n=U(e,n,s,t);U(e,n,null,t)}else n==null||n===""||!e.firstChild?e.appendChild(t):e.replaceChild(t,e.firstChild);n=t}}return n}function de(e,t,n,s){let r=!1;for(let i=0,o=t.length;i<o;i++){let l=t[i],u=n&&n[i],f;if(!(l==null||l===!0||l===!1))if((f=typeof l)=="object"&&l.nodeType)e.push(l);else if(Array.isArray(l))r=de(e,l,u)||r;else if(f==="function")if(s){for(;typeof l=="function";)l=l();r=de(e,Array.isArray(l)?l:[l],Array.isArray(u)?u:[u])||r}else e.push(l),r=!0;else{const c=String(l);u&&u.nodeType===3&&u.data===c?e.push(u):e.push(document.createTextNode(c))}}return r}function Ae(e,t,n=null){for(let s=0,r=t.length;s<r;s++)e.insertBefore(t[s],n)}function U(e,t,n,s){if(n===void 0)return e.textContent="";const r=s||document.createTextNode("");if(t.length){let i=!1;for(let o=t.length-1;o>=0;o--){const l=t[o];if(r!==l){const u=l.parentNode===e;!i&&!o?u?e.replaceChild(r,l):e.insertBefore(r,n):u&&l.remove()}else i=!0}}else e.insertBefore(r,n);return[r]}const _t="http://www.w3.org/2000/svg";function vt(e,t=!1){return t?document.createElementNS(_t,e):document.createElement(e)}function cn(e){const{useShadow:t}=e,n=document.createTextNode(""),s=()=>e.mount||document.body,r=ye();let i,o=!!m.context;return he(()=>{o&&(ye().user=o=!1),i||(i=Je(r,()=>I(()=>e.children)));const l=s();if(l instanceof HTMLHeadElement){const[u,f]=O(!1),c=()=>f(!0);F(d=>b(l,()=>u()?d():i(),null)),q(c)}else{const u=vt(e.isSVG?"g":"div",e.isSVG),f=t&&u.attachShadow?u.attachShadow({mode:"open"}):u;Object.defineProperty(u,"_$host",{get(){return n.parentNode},configurable:!0}),b(f,i),l.appendChild(u),e.ref&&e.ref(u),q(()=>l.removeChild(u))}},void 0,{render:!o}),n}const kt="modulepreload",Nt=function(e){return"/"+e},Se={},Pt=function(t,n,s){let r=Promise.resolve();if(n&&n.length>0){const i=document.getElementsByTagName("link");r=Promise.all(n.map(o=>{if(o=Nt(o),o in Se)return;Se[o]=!0;const l=o.endsWith(".css"),u=l?'[rel="stylesheet"]':"";if(!!s)for(let d=i.length-1;d>=0;d--){const a=i[d];if(a.href===o&&(!l||a.rel==="stylesheet"))return}else if(document.querySelector(`link[href="${o}"]${u}`))return;const c=document.createElement("link");if(c.rel=l?"stylesheet":kt,l||(c.as="script",c.crossOrigin=""),c.href=o,document.head.appendChild(c),l)return new Promise((d,a)=>{c.addEventListener("load",d),c.addEventListener("error",()=>a(new Error(`Unable to preload CSS for ${o}`)))})}))}return r.then(()=>t()).catch(i=>{const o=new Event("vite:preloadError",{cancelable:!0});if(o.payload=i,window.dispatchEvent(o),!o.defaultPrevented)throw i})},Tt=A('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="m8.382 17.025l-1.407-1.4L10.593 12L6.975 8.4L8.382 7L12 10.615L15.593 7L17 8.4L13.382 12L17 15.625l-1.407 1.4L12 13.41z">'),Ot=(e={})=>(()=>{const t=Tt();return $t(t,e,!0,!0),t})(),It=A("<button rounded-md p-2 bg-zinc-700 hover:bg-zinc-600>"),Mt=A("<div flex items-center flex-wrap gap-4 justify-between mb-2 p-4 md:px-6><h1 m-0>Inbox</h1><div flex justify-center items-center gap-2><input outline-none rounded-md bg-zinc-800 hover:bg-zinc-700 transition-colors font-normal placeholder=Search sm:w-auto px-4 py-2>");function jt({setSearch:e,search:t}){return(()=>{const n=Mt(),s=n.firstChild,r=s.nextSibling,i=r.firstChild;return b(r,h(G,{get when(){return t()},get children(){const o=It();return o.$$click=()=>e(""),b(o,h(Ot,{})),o}}),i),i.$$input=o=>e(o.target.value),i.style.setProperty("width","200px"),j(()=>i.value=t()),n})()}ge(["click","input"]);const zt=A("<span font-normal text-zinc-500>"),Rt=A("<span font-bold text-zinc-200>");function Dt({address:e}){const t=()=>e().name,n=()=>e().address;return(()=>{const s=Rt();return b(s,t,null),b(s,h(G,{get when(){return n()},get children(){const r=zt();return b(r,()=>` <${n()}>`),r}}),null),s})()}const Ut=A("<div flex items-center><p m-0 truncate w=[300px]></p><p my-0 mx-3 truncate flex-1>"),Vt=A("<div flex items-center><p m-0 truncate flex-1>"),Ce=A("<p m-0 truncate>"),Bt=A("<p m-0 mb-1 truncate>"),Ht=A("<button class=block cursor-pointer w-full text-left border-0 border-t border-t-solid border-zinc-800 py-3 px-3 md:px-6 rounded-none text-zinc-500>"),Ft=A("<span text-zinc-500>Unknown address"),Wt=A("<span font-bold text-zinc-200>"),qt=A("<span font-normal text-zinc-400> - "),Kt=A("<p m-0 text-zinc-500 text=13px>");function Gt({display:e,selected:t,email:n,onClick:s}){const r=()=>n.from[0],i=I(()=>{const o=new Date(n.realDate);return o.toLocaleDateString()+" "+o.toLocaleTimeString()});return(()=>{const o=Ht();return Ue(o,"click",s,!0),b(o,h(De,{get children(){return[h(W,{get when(){return e()==="lg"},get children(){const l=Ut(),u=l.firstChild,f=u.nextSibling;return b(u,h(ce,{from:r})),b(f,h(ue,{email:()=>n})),b(l,h(fe,{date:i}),null),l}}),h(W,{get when(){return e()==="md"},get children(){return[(()=>{const l=Vt(),u=l.firstChild;return b(u,h(ce,{from:r})),b(l,h(fe,{date:i}),null),l})(),(()=>{const l=Ce();return b(l,h(ue,{email:()=>n})),l})()]}}),h(W,{get when(){return e()==="sm"},get children(){return[(()=>{const l=Ce();return b(l,h(ce,{from:r})),l})(),(()=>{const l=Bt();return b(l,h(ue,{email:()=>n})),l})(),h(fe,{date:i})]}})]}})),j(()=>K(o,"bg",t()?"zinc-800":"zinc-950")),o})()}function ce({from:e}){return h(G,{get when(){return e()},get fallback(){return Ft()},children:t=>h(Dt,{address:t})})}function ue({email:e}){return[(()=>{const t=Wt();return b(t,()=>e().subject),t})(),h(G,{get when(){return e().bodyHint},get children(){const t=qt();return t.firstChild,b(t,()=>e().bodyHint,null),t}})]}function fe({date:e,fullWidth:t}){return(()=>{const n=Kt();return K(n,"w",t?"130px":"auto"),b(n,e),n})()}ge(["click"]);const Xt=A("<div>"),Qt=A("<p p-3 text-zinc-400>Loading..."),Yt=A("<p p-3 text-zinc-400>No emails");function Ve({emails:e,loading:t}){const[n,s]=pe(Fe),[r,i]=pe(We),[o,l]=O(),u=a=>{if(a)return a>900?"lg":a>600?"md":"sm"};let f=null;const c=()=>{const a=f==null?void 0:f.getBoundingClientRect().width;a&&l(u(a))},d=({target:a,key:y})=>{const E=n(),g=e();if(!E||!g||g.length<2)return;const N=y==="ArrowUp"?-1:y==="ArrowDown"?1:0;if(N!==0){if(a){const k=a;if(new Set(["input","textarea","area","select"]).has(k.tagName.toLowerCase()))return}for(let k=0;k<g.length;k++)if(g[k].id===E.id){const L=g[k+N];L&&s.select(L);break}}};return Ne(()=>{window.addEventListener("resize",c),window.addEventListener("keydown",d),c()}),q(()=>{window.removeEventListener("resize",c),window.removeEventListener("keydown",d)}),(()=>{const a=Xt();return At(y=>f=y,a),b(a,h(jt,{search:r,get setSearch(){return i.set}}),null),b(a,h(G,{get when(){return!t||!t()},get fallback(){return Qt()},get children(){return h(ot,{get each(){return e()},get fallback(){return Yt()},children:y=>{const E=()=>{var g;return y.id==((g=n())==null?void 0:g.id)};return h(Gt,{display:()=>o(),selected:E,email:y,onClick:()=>E()?s.deSelect():s.select(y)})}})}}),null),a})()}function Be(e,t=!1){return t?(location.protocol==="https:"?"wss://":"ws://")+location.host+e:e}function Le(e,t){return window.fetch(Be(e),t)}class Jt{constructor(t){X(this,"closed",!1);X(this,"socket",null);X(this,"newEmail");this.newEmail=t}async start(){if(this.socket)throw"cannot start a Websocket twice";let t=!0;for(;await new Promise(n=>{this.socket=new WebSocket(Be("/api/emails-events",!0)),this.socket.onopen=()=>{this.closed&&!t&&this.newEmail()},this.socket.onmessage=s=>{!this.closed&&s.data==="new-email"&&this.newEmail()},this.socket.onclose=()=>setTimeout(n,5e3)}),!this.closed;)t=!1}close(){var t;this.closed=!0,(t=this.socket)==null||t.close()}}const Zt=A("<div h-full w-full overflow-hidden>"),en=A("<div flex items-stretch><div h-screen self-stretch overflow-y-auto w=[400px] border-0 border-r border-r-solid border-zinc-700>"),He=lt(()=>Pt(()=>import("./Show-pm1-B8_e.js"),__vite__mapDeps([])).then(e=>({default:e.ShowEmail}))),Fe=Pe([]),We=Pe([]);function tn(e){let t=!1,n=!1;return async()=>{if(t){n=!0;return}for(t=!0,n=!0;n;){n=!1;try{await e()}catch(s){console.error(s)}await new Promise(s=>setTimeout(s,500))}t=!1}}const _e=()=>window.innerWidth<840;function nn(){const[e,t]=O(),[n,s]=O(),[r,i]=O(""),[o,l]=O(_e()),u=tn(async()=>{var C;let g="/api/emails";(C=r())!=null&&C.length&&(g+="?search="+encodeURIComponent(r()));const k=await(await Le(g)).json();t(k)}),f=new Jt(u);he(()=>{r(),u()});const c=()=>{const g=_e();g!==o()&&l(g)};Ne(()=>{f.start(),window.addEventListener("resize",c)}),q(()=>{f.close(),window.removeEventListener("resize",c)});const d=async()=>{var N;const g=(N=n())==null?void 0:N.id;g&&(s(void 0),await Le("/api/emails/"+g,{method:"DELETE"}),await u())},a=()=>[n,{delete:d,select:s,deSelect:()=>s(void 0)}],y=()=>[r,{set:i}],E=()=>e()!==void 0&&n()!==void 0;return h(Fe.Provider,{get value(){return a()},get children(){return h(We.Provider,{get value(){return y()},get children(){const g=Zt();return b(g,h(De,{get fallback(){return h(Ve,{emails:()=>e()??[],loading:()=>!e()})},get children(){return[h(W,{get when(){return I(()=>!!E())()&&o()},get children(){return h(He,{})}}),h(W,{get when(){return E()},get children(){return h(sn,{emails:e,selectedEmail:()=>n()})}})]}})),g}})}})}function sn({emails:e}){return(()=>{const t=en(),n=t.firstChild;return b(n,h(Ve,{emails:e})),b(t,h(He,{}),null),t})()}const rn=document.getElementById("root");yt(()=>h(nn,{}),rn);export{Dt as E,ot as F,Ot as I,W as M,cn as P,Fe as S,Pt as _,h as a,Ue as b,I as c,ge as d,De as e,q as f,O as g,he as h,b as i,j,K as k,G as l,At as m,Be as n,Ne as o,lt as p,on as q,Et as r,$t as s,A as t,pe as u,Le as v};
function __vite__mapDeps(indexes) {
  if (!__vite__mapDeps.viteFileDeps) {
    __vite__mapDeps.viteFileDeps = []
  }
  return indexes.map((i) => __vite__mapDeps.viteFileDeps[i])
}