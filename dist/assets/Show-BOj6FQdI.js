const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/Code-Bxr2gThp.js","assets/index-8HhG6-LH.js","assets/index-Bv90_tQ-.css","assets/Code-DV6HP0dz.css"])))=>i.map(i=>d[i]);
import{d as x,u as A,S,c as V,a as r,E,b as $,i as d,I as L,F as P,t as s,s as k,M as p,e as y,o as j,f as F,P as D,g as T,h as C,j as g,k as q,l as w,m as R,n as B,p as U,_ as O,q as N}from"./index-8HhG6-LH.js";var W=s("<div bg-zinc-900><div flex justify-between items-center p-4 border-0 border-b border-b-solid border-zinc-800><h2 font-bold m-0 flex items-center><button p-0 text-zinc-400 inline-flex mt-1 mr-2></button></h2><div><button bg-zinc-700 hover:bg-zinc-600>Delete</button></div></div><div border-0 border-b border-b-solid border-zinc-800 p-4 flex flex-wrap>"),G=s("<span mr-3>: ");function J({email:e,emailRemainder:t}){const[o,i]=A(S),c=V(()=>{const n=e(),a=[],l=(m,h)=>{h&&h.length>0&&a.push({name:m,value:h.map(v=>r(E,{address:()=>v}))})};n.sender&&a.push({name:"Sender",value:[r(E,{address:()=>n.sender})]}),l("From",n.from),l("To",n.to);const u=t();return u&&(l("ReplyTo",u.replyTo),l("CC",u.cc),l("BCC",u.bcc)),a});return(()=>{var n=W(),a=n.firstChild,l=a.firstChild,u=l.firstChild,m=l.nextSibling,h=m.firstChild,v=a.nextSibling;return $(u,"click",i.deSelect,!0),d(u,r(L,{})),d(l,()=>e().subject,null),$(h,"click",i.delete,!0),d(v,r(P,{get each(){return c()},children:b=>(()=>{var f=G(),_=f.firstChild;return d(f,()=>b.name,_),d(f,()=>b.value,null),f})()})),n})()}x(["click"]);var K=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M5 21q-.825 0-1.412-.587T3 19V5q0-.825.588-1.412T5 3h14q.825 0 1.413.588T21 5v14q0 .825-.587 1.413T19 21zm0-2h14V5H5zm1-2h12l-3.75-5l-3 4L9 13zm-1 2V5z">');const Q=(e={})=>(()=>{var t=K();return k(t,e,!0),t})();var X=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M6 22q-.825 0-1.412-.587T4 20V4q0-.825.588-1.412T6 2h9l5 5v13q0 .825-.587 1.413T18 22zm0-2h12V8h-4V4H6zm6-1q1.675 0 2.838-1.175T16 15v-4h-2v4q0 .825-.575 1.413T12 17q-.825 0-1.412-.587T10 15V9.5q0-.225.15-.363T10.5 9q.225 0 .363.138T11 9.5V15h2V9.5q0-1.05-.725-1.775T10.5 7t-1.775.725T8 9.5V15q0 1.65 1.175 2.825T12 19M6 4v4zv16z">');const Y=(e={})=>(()=>{var t=X();return k(t,e,!0),t})();var Z=s('<svg viewBox="0 0 24 24"width=1.2em height=1.2em><path fill=currentColor d="M9 12.5h1v-2h1q.425 0 .713-.288T12 9.5v-1q0-.425-.288-.712T11 7.5H9zm1-3v-1h1v1zm3 3h2q.425 0 .713-.288T16 11.5v-3q0-.425-.288-.712T15 7.5h-2zm1-1v-3h1v3zm3 1h1v-2h1v-1h-1v-1h1v-1h-2zM8 18q-.825 0-1.412-.587T6 16V4q0-.825.588-1.412T8 2h12q.825 0 1.413.588T22 4v12q0 .825-.587 1.413T20 18zm0-2h12V4H8zm-4 6q-.825 0-1.412-.587T2 20V6h2v14h14v2zM8 4v12z">');const ee=(e={})=>(()=>{var t=Z();return k(t,e,!0),t})();function M({contentType:e}){const t={height:"30px",width:"30px"};return r(y,{get fallback(){return r(Y,{style:t})},get children(){return[r(p,{get when(){return I(e())},get children(){return r(Q,{style:t})}}),r(p,{get when(){return H(e())},get children(){return r(ee,{style:t})}})]}})}var te=s("<button p-4 border-0 w=180px><div truncate>");const H=e=>e==="application/pdf"||e==="application/x-pdf"||e==="application/x-bzpdf",I=e=>e.startsWith("image/");function re({attachment:e,onclick:t}){return(()=>{var o=te(),i=o.firstChild;return $(o,"click",t,!0),d(o,r(M,{contentType:()=>e().contentType}),i),d(i,()=>e().filename),o})()}x(["click"]);var ne=s("<div fixed z-100 top-0 left-0 h-screen w-screen p-10 bg=black/40><div mx-auto bg-zinc-900 rounded max-w=800px h-full shadow-xl overflow-hidden flex flex-col><div p-6 border-0 border-zinc-800 border-b border-b-solid><p m-0></p></div><div py-4 px-6 flex justify-end gap-2 items-center border-0 border-zinc-800 border-t border-t-solid><button bg-zinc-700 hover:bg-zinc-600>Close</button><button bg-zinc-700 hover:bg-zinc-600>Download"),le=s("<div h-full w-full bg-center bg-contain bg-no-repeat>"),ie=s("<iframe h-full>"),oe=s("<p flex flex-col justify-center items-center h-full text-zinc-400><span mt-3>Content not displayable");function ce({onClose:e,attachmentUrl:t,attachment:o}){let i;const c=l=>{var u;(!i||i===l.target||i===((u=l.target.children)==null?void 0:u[0]))&&e()},n=()=>{const l=document.createElement("a");l.download=o().filename,l.href=t(),l.click()},a=l=>{l.key==="Escape"&&e()};return j(()=>{document.addEventListener("keydown",a)}),F(()=>{document.removeEventListener("keydown",a)}),r(D,{get mount(){return document.body},get children(){var l=ne(),u=l.firstChild,m=u.firstChild,h=m.firstChild,v=m.nextSibling,b=v.firstChild,f=b.nextSibling;l.$$click=c;var _=i;return typeof _=="function"?R(_,l):i=l,d(h,()=>o().filename),d(u,r(ae,{attachment:o,attachmentUrl:t}),v),$(b,"click",e,!0),f.$$click=n,l}})}function ae({attachment:e,attachmentUrl:t}){const[o,i]=T(void 0),c=async()=>{const a=await(await fetch(t())).blob(),l=URL.createObjectURL(a);i(l)};return C(()=>{c()}),r(w,{get when(){return o()},get children(){return r(y,{get fallback(){return r(de,{contentType:()=>e().contentType})},get children(){return[r(p,{get when(){return I(e().contentType)},get children(){var n=le();return g(a=>(a=`url(${o()})`)!=null?n.style.setProperty("background-image",a):n.style.removeProperty("background-image")),n}}),r(p,{get when(){return H(e().contentType)},get children(){var n=ie();return g(()=>q(n,"src",o())),n}})]}})}})}function de({contentType:e}){return(()=>{var t=oe(),o=t.firstChild;return d(t,r(M,{contentType:e}),o),t})()}x(["click"]);var se=s("<div border-0 border-b border-b-solid border-zinc-800 p-4><p m-0 text-zinc-500>Attachments:</p><div flex flex-wrap gap-2 mt-1>");function ue({emailRemainder:e,email:t}){const[o,i]=T(),c=()=>{var n,a;return(((a=(n=e())==null?void 0:n.attachments)==null?void 0:a.length)??0)>0};return C(()=>{t(),i(void 0)}),[r(w,{get when(){return c()},get children(){var n=se(),a=n.firstChild,l=a.nextSibling;return d(l,r(P,{get each(){return e().attachments},children:(u,m)=>r(re,{attachment:()=>u,onclick:()=>i(m)})})),n}}),r(w,{get when(){return V(()=>o()!==void 0)()&&c()},get children(){return r(ce,{onClose:()=>i(void 0),attachmentUrl:()=>B("/api/emails/"+t().id+"/attachments/"+o()),attachment:()=>e().attachments[o()]})}})]}var he=s("<div p-4 flex-1 overflow-y-auto><pre>"),me=s("<div flex-1>"),pe=s("<button border-0>"),ve=s("<iframe flex-1 w-full bg-white h-100>"),fe=s("<pre p-4 pt-0 overflow-y-auto>"),ge=s("<div p-4 pt-0 overflow-y-auto>"),be=s("<div flex flex-col h-full><div flex gap-1 p-4>"),$e=s("<p p-4 pt-0 text-zinc-500>This email does not have a plain text fallback 😞"),_e=s("<p italic text-zinc-500>");const we=U(()=>O(()=>import("./Code-Bxr2gThp.js"),__vite__mapDeps([0,1,2,3])).then(e=>({default:e.Code})));function xe({email:e,emailRemainder:t}){const o=()=>{var c;return((c=t())==null?void 0:c.htmlBody)??""},i=()=>{var c;return((c=t())==null?void 0:c.textBody)??""};return r(y,{get children(){return[r(p,{get when(){return e().bodyType==="text"},get children(){var c=he(),n=c.firstChild;return d(c,r(Te,{kind:"text/plain"}),null),g(()=>n.innerText=i()),c}}),r(p,{get when(){return e().bodyType==="html"},get children(){var c=me();return d(c,r(ye,{email:e,html:o,plain:i})),c}})]}})}function z({onclick:e,selected:t,children:o}){return(()=>{var i=pe();return $(i,"click",e,!0),d(i,o),g(()=>q(i,"text",t()?"indigo-400":"zinc-200")),i})()}function ye({email:e,html:t,plain:o}){const[i,c]=T("preview"),n=()=>i()==="preview",a=()=>i()==="plain",l=()=>i()==="html";return(()=>{var u=be(),m=u.firstChild;return d(m,r(z,{onclick:()=>c("preview"),selected:n,children:"Preview"}),null),d(m,r(z,{onclick:()=>c("plain"),selected:a,children:"Plain text"}),null),d(m,r(z,{onclick:()=>c("html"),selected:l,children:"HTML"}),null),d(u,r(y,{get children(){return[r(p,{get when(){return n()},get children(){var h=ve();return g(()=>q(h,"src",B(`/api/emails/${e().id}/page`))),h}}),r(p,{get when(){return a()},get children(){return r(w,{get when(){return o()},get fallback(){return $e()},get children(){var h=fe();return g(()=>h.innerText=o()),h}})}}),r(p,{get when(){return l()},get children(){var h=ge();return d(h,r(we,{lang:"xml",code:t})),h}})]}}),null),u})()}function Te({kind:e}){return(()=>{var t=_e();return d(t,e),t})()}x(["click"]);var ze=s("<div flex flex-col overflow-y-hidden h-screen>");function Ce({}){const[e]=A(S),t=()=>e(),[o,i]=T(),c=async n=>{i(void 0);const a=await N(`/api/emails/${n}/remainder`),l=await a.json();if(a.status!==200)throw l.error;i(l)};return C(()=>{const n=e();n&&c(n.id)}),(()=>{var n=ze();return d(n,r(J,{email:t,emailRemainder:o}),null),d(n,r(ue,{email:t,emailRemainder:o}),null),d(n,r(xe,{email:t,emailRemainder:o}),null),n})()}export{Ce as Show};