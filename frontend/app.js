const el = (id)=>document.getElementById(id)
const sessionsEl = el('sessions'), panel = el('chatPanel'), input = el('input'), send = el('send'), newChat = el('newChat')
const modelSelect = el('modelSelect'), streamToggle = el('streamToggle')

let sessions = JSON.parse(localStorage.getItem('sql_sessions')||'[]')
let current = localStorage.getItem('sql_current') || null

function save(){ localStorage.setItem('sql_sessions', JSON.stringify(sessions)); localStorage.setItem('sql_current', current) }
function renderSessions(){
  sessionsEl.innerHTML = ''
  sessions.forEach((s, i)=>{
    const d=document.createElement('div'); d.className='session'; d.textContent=s.title||('Session '+(i+1))
    d.onclick=()=>{ current=i; renderMessages(); save() }
    sessionsEl.appendChild(d)
  })
}
function newSession(){
  const s={title:'New Chat', messages:[]}
  sessions.unshift(s); current=0; save(); renderSessions(); renderMessages()
}
function renderMessages(){
  panel.innerHTML=''
  if(current==null){ return }
  const msgs=sessions[current].messages
  msgs.forEach(m=>{ panel.appendChild(renderMessage(m)) })
  panel.scrollTop = panel.scrollHeight
}
function renderMessage(m){
  const d=document.createElement('div'); d.className='msg '+(m.role==='user'?'user':'assistant')
  if(m.role==='assistant' && /;|\n/.test(m.content)){
    const pre=document.createElement('div'); pre.className='codeblock'; pre.textContent=m.content
    const btn=document.createElement('button'); btn.className='copyBtn'; btn.textContent='Copy'
    btn.onclick=()=>navigator.clipboard.writeText(m.content)
    pre.appendChild(btn); d.appendChild(pre)
  } else {
    d.textContent = m.content
  }
  return d
}

async function sendMessage(){
  const text = input.value.trim(); if(!text) return
  if(current==null) newSession()
  const roleUser={role:'user',content:text}
  sessions[current].messages.push(roleUser); renderSessions(); renderMessages(); input.value=''
  const model = modelSelect.value
  const streamOn = streamToggle.dataset.on === 'true'
  if(streamOn){
    await streamRequest(model, text)
  } else {
    await fetchRequest(model, text)
  }
  save()
}

async function fetchRequest(model, prompt){
  const res = await fetch('/api/chat',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({model,prompt,stream:false})})
  const j = await res.json()
  const assistant = {role:'assistant',content:j.text}
  sessions[current].messages.push(assistant); renderMessages()
}

async function streamRequest(model, prompt){
  const url = `/api/chat/stream?model=${encodeURIComponent(model)}&prompt=${encodeURIComponent(prompt)}`
  const es = new EventSource(url)
  if(current==null) newSession()
  const assistant={role:'assistant',content:''}
  sessions[current].messages.push(assistant); renderMessages()
  es.onmessage = (e)=>{
    const chunk = e.data
    assistant.content += chunk
    renderMessages()
  }
  await new Promise((res, rej)=>{
    es.onerror = ()=>{ es.close(); res() }
    // close after 1s idle if needed
    setTimeout(()=>{ es.close(); res() }, 120*1000)
  })
}

send.addEventListener('click', sendMessage)
input.addEventListener('keydown', (e)=>{ if(e.key==='Enter' && !e.shiftKey){ e.preventDefault(); sendMessage() } })
newChat.addEventListener('click', newSession)
streamToggle.dataset.on = 'true'
streamToggle.addEventListener('click', ()=>{
  const on = streamToggle.dataset.on === 'true'
  streamToggle.dataset.on = (!on) + ''
  streamToggle.textContent = 'Stream: ' + (on? 'OFF' : 'ON')
})

renderSessions()
renderMessages()
