const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/Code-DenXeoXT.js","assets/index-6eJkkb5r.js","assets/index-BfofDsXp.css","assets/Code-DV6HP0dz.css"])))=>i.map(i=>d[i]);
import{d as b,u as A,S as V,c as S,a as r,E,i as d,b as _,F as P,t as s,s as z,M as p,e as w,o as j,f as I,P as F,g as x,h as k,j as g,k as C,l as $,m as D,n as B,p as R,_ as U,q as O}from"./index-6eJkkb5r.js";var N=s("<div bg-zinc-900><div flex justify-between items-center p-4 border-0 border-b border-b-solid border-zinc-800><h2 font-bold m-0></h2><div><button border border-zinc-400>Delete</button></div></div><div border-0 border-b border-b-solid border-zinc-800 p-4 flex flex-wrap>"),W=s("<span mr-3>: ");function G({email:e,emailRemainder:t}){const[o,i]=A(V),c=S(()=>{const n=e(),a=[],l=(m,h)=>{h&&h.length>0&&a.push({name:m,value:h.map(v=>r(E,{address:()=>v}))})};n.sender&&a.push({name:"Sender",value:[r(E,{address:()=>n.sender})]}),l("From",n.from),l("To",n.to);const u=t();return u&&(l("ReplyTo",u.replyTo),l("CC",u.cc),l("BCC",u.bcc)),a});return(()=>{var n=N(),a=n.firstChild,l=a.firstChild,u=l.nextSibling,m=u.firstChild,h=a.nextSibling;return d(l,()=>e().subject),_(m,"click",i.delete,!0),d(h,r(P,{get each(){return c()},children:v=>(()=>{var f=W(),y=f.firstChild;return d(f,()=>v.name,y),d(f,()=>v.value,null),f})()})),n})()}b(["click"]);var J=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M5 21q-.825 0-1.412-.587T3 19V5q0-.825.588-1.412T5 3h14q.825 0 1.413.588T21 5v14q0 .825-.587 1.413T19 21zm0-2h14V5H5zm1-2h12l-3.75-5l-3 4L9 13zm-1 2V5z">');const K=(e={})=>(()=>{var t=J();return z(t,e,!0),t})();var Q=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M6 22q-.825 0-1.412-.587T4 20V4q0-.825.588-1.412T6 2h9l5 5v13q0 .825-.587 1.413T18 22zm0-2h12V8h-4V4H6zm6-1q1.675 0 2.838-1.175T16 15v-4h-2v4q0 .825-.575 1.413T12 17q-.825 0-1.412-.587T10 15V9.5q0-.225.15-.363T10.5 9q.225 0 .363.138T11 9.5V15h2V9.5q0-1.05-.725-1.775T10.5 7t-1.775.725T8 9.5V15q0 1.65 1.175 2.825T12 19M6 4v4zv16z">');const X=(e={})=>(()=>{var t=Q();return z(t,e,!0),t})();var Y=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M9 12.5h1v-2h1q.425 0 .713-.288T12 9.5v-1q0-.425-.288-.712T11 7.5H9zm1-3v-1h1v1zm3 3h2q.425 0 .713-.288T16 11.5v-3q0-.425-.288-.712T15 7.5h-2zm1-1v-3h1v3zm3 1h1v-2h1v-1h-1v-1h1v-1h-2zM8 18q-.825 0-1.412-.587T6 16V4q0-.825.588-1.412T8 2h12q.825 0 1.413.588T22 4v12q0 .825-.587 1.413T20 18zm0-2h12V4H8zm-4 6q-.825 0-1.412-.587T2 20V6h2v14h14v2zM8 4v12z">');const Z=(e={})=>(()=>{var t=Y();return z(t,e,!0),t})();function M({contentType:e}){const t={height:"30px",width:"30px"};return r(w,{get fallback(){return r(X,{style:t})},get children(){return[r(p,{get when(){return L(e())},get children(){return r(K,{style:t})}}),r(p,{get when(){return H(e())},get children(){return r(Z,{style:t})}})]}})}var ee=s("<button p-4 border-0 w=180px><div truncate>");const H=e=>e==="application/pdf"||e==="application/x-pdf"||e==="application/x-bzpdf",L=e=>e.startsWith("image/");function te({attachment:e,onclick:t}){return(()=>{var o=ee(),i=o.firstChild;return _(o,"click",t,!0),d(o,r(M,{contentType:()=>e().contentType}),i),d(i,()=>e().filename),o})()}b(["click"]);var re=s("<div fixed z-100 top-0 left-0 h-screen w-screen p-10 bg=black/40><div mx-auto bg-zinc-900 rounded max-w=800px h-full shadow-xl overflow-hidden flex flex-col><div p-6 border-0 border-zinc-800 border-b border-b-solid><p m-0></p></div><div py-4 px-6 flex justify-end gap-2 items-center border-0 border-zinc-800 border-t border-t-solid><button border-zinc-500>Close</button><button border-zinc-500>Download"),ne=s("<div h-full w-full bg-center bg-contain bg-no-repeat>"),le=s("<iframe h-full>"),ie=s("<p flex flex-col justify-center items-center h-full text-zinc-400><span mt-3>Content not displayable");function oe({onClose:e,attachmentUrl:t,attachment:o}){let i;const c=l=>{var u;(!i||i===l.target||i===((u=l.target.children)==null?void 0:u[0]))&&e()},n=()=>{const l=document.createElement("a");l.download=o().filename,l.href=t(),l.click()},a=l=>{l.key==="Escape"&&e()};return j(()=>{document.addEventListener("keydown",a)}),I(()=>{document.removeEventListener("keydown",a)}),r(F,{get mount(){return document.body},get children(){var l=re(),u=l.firstChild,m=u.firstChild,h=m.firstChild,v=m.nextSibling,f=v.firstChild,y=f.nextSibling;l.$$click=c;var q=i;return typeof q=="function"?D(q,l):i=l,d(h,()=>o().filename),d(u,r(ce,{attachment:o,attachmentUrl:t}),v),_(f,"click",e,!0),y.$$click=n,l}})}function ce({attachment:e,attachmentUrl:t}){const[o,i]=x(void 0),c=async()=>{const a=await(await fetch(t())).blob(),l=URL.createObjectURL(a);i(l)};return k(()=>{c()}),r($,{get when(){return o()},get children(){return r(w,{get fallback(){return r(ae,{contentType:()=>e().contentType})},get children(){return[r(p,{get when(){return L(e().contentType)},get children(){var n=ne();return g(a=>(a=`url(${o()})`)!=null?n.style.setProperty("background-image",a):n.style.removeProperty("background-image")),n}}),r(p,{get when(){return H(e().contentType)},get children(){var n=le();return g(()=>C(n,"src",o())),n}})]}})}})}function ae({contentType:e}){return(()=>{var t=ie(),o=t.firstChild;return d(t,r(M,{contentType:e}),o),t})()}b(["click"]);var de=s("<div border-0 border-b border-b-solid border-zinc-800 p-4><p m-0 text-zinc-500>Attachments:</p><div flex flex-wrap gap-2 mt-1>");function se({emailRemainder:e,email:t}){const[o,i]=x(),c=()=>{var n,a;return(((a=(n=e())==null?void 0:n.attachments)==null?void 0:a.length)??0)>0};return k(()=>{t(),i(void 0)}),[r($,{get when(){return c()},get children(){var n=de(),a=n.firstChild,l=a.nextSibling;return d(l,r(P,{get each(){return e().attachments},children:(u,m)=>r(te,{attachment:()=>u,onclick:()=>i(m)})})),n}}),r($,{get when(){return S(()=>o()!==void 0)()&&c()},get children(){return r(oe,{onClose:()=>i(void 0),attachmentUrl:()=>B("/api/emails/"+t().id+"/attachments/"+o()),attachment:()=>e().attachments[o()]})}})]}var ue=s("<div p-4 flex-1 overflow-y-auto><pre>"),he=s("<div flex-1>"),me=s("<button border-0>"),pe=s("<iframe flex-1 w-full bg-white h-100>"),ve=s("<pre p-4 pt-0 overflow-y-auto>"),fe=s("<div p-4 pt-0 overflow-y-auto>"),ge=s("<div flex flex-col h-full><div flex gap-1 p-4>"),$e=s("<p p-4 pt-0 text-zinc-500>This email does not have a plain text fallback 😞"),be=s("<p italic text-zinc-500>");const _e=R(()=>U(()=>import("./Code-DenXeoXT.js"),__vite__mapDeps([0,1,2,3])).then(e=>({default:e.Code})));function we({email:e,emailRemainder:t}){const o=()=>{var c;return((c=t())==null?void 0:c.htmlBody)??""},i=()=>{var c;return((c=t())==null?void 0:c.textBody)??""};return r(w,{get children(){return[r(p,{get when(){return e().bodyType==="text"},get children(){var c=ue(),n=c.firstChild;return d(c,r(ye,{kind:"text/plain"}),null),g(()=>n.innerText=i()),c}}),r(p,{get when(){return e().bodyType==="html"},get children(){var c=he();return d(c,r(xe,{email:e,html:o,plain:i})),c}})]}})}function T({onclick:e,selected:t,children:o}){return(()=>{var i=me();return _(i,"click",e,!0),d(i,o),g(()=>C(i,"text",t()?"indigo-400":"zinc-200")),i})()}function xe({email:e,html:t,plain:o}){const[i,c]=x("preview"),n=()=>i()==="preview",a=()=>i()==="plain",l=()=>i()==="html";return(()=>{var u=ge(),m=u.firstChild;return d(m,r(T,{onclick:()=>c("preview"),selected:n,children:"Preview"}),null),d(m,r(T,{onclick:()=>c("plain"),selected:a,children:"Plain text"}),null),d(m,r(T,{onclick:()=>c("html"),selected:l,children:"HTML"}),null),d(u,r(w,{get children(){return[r(p,{get when(){return n()},get children(){var h=pe();return g(()=>C(h,"src",B(`/api/emails/${e().id}/page`))),h}}),r(p,{get when(){return a()},get children(){return r($,{get when(){return o()},get fallback(){return $e()},get children(){var h=ve();return g(()=>h.innerText=o()),h}})}}),r(p,{get when(){return l()},get children(){var h=fe();return d(h,r(_e,{lang:"xml",code:t})),h}})]}}),null),u})()}function ye({kind:e}){return(()=>{var t=be();return d(t,e),t})()}b(["click"]);var Te=s("<div flex flex-col overflow-y-hidden h-screen>");function ke({}){const[e]=A(V),t=()=>e(),[o,i]=x(),c=async n=>{i(void 0);const a=await O(`/api/emails/${n}/remainder`),l=await a.json();if(a.status!==200)throw l.error;i(l)};return k(()=>{const n=e();n&&c(n.id)}),(()=>{var n=Te();return d(n,r(G,{email:t,emailRemainder:o}),null),d(n,r(se,{email:t,emailRemainder:o}),null),d(n,r(we,{email:t,emailRemainder:o}),null),n})()}export{ke as Show};