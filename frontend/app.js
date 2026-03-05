let ws = new WebSocket("ws://localhost:8002/ws")

ws.onmessage = (event)=>{

let div=document.createElement("div")
div.className="message"
div.innerText=event.data

document.getElementById("messages").appendChild(div)

}

function send(){

let msg=document.getElementById("msg").value

if(msg==="") return

let div=document.createElement("div")
div.className="message self"
div.innerText=msg

document.getElementById("messages").appendChild(div)

ws.send(msg)

document.getElementById("msg").value=""

}