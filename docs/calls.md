# Ligações

## Ligações Iniciadas pelo Cliente

```js
// ws = new WebSocket("path-to-wabaman-wss-endpoint");

// eventos relevantes:
const MSG_INCOMING_CALL_FROM_CLIENT = 30;
const MSG_SETUP_CALL_FROM_BROWSER = 31;
const MSG_TERMINATE_CALL = 32;
const MSG_CALL_CONSUMED = 33;
const MSG_ACCEPT_CALL = 34;
const MSG_REJECT_CALL = 35;
const MSG_SEND_BROWSER_CANDIDATE = 36;
const MSG_CALL_STARTED = 37;
const MSG_CALL_ENDED = 38;
const MSG_CALL_ON_ANSWER_SDP = 39;
const MSG_CALL_START_TIMER = 40;
```

### Recebendo uma Ligação

```js
ws.onmessage = function (event) {
    const data = JSON.parse(event.data);
    
    if (data.type === MSG_INCOMING_CALL_FROM_CLIENT) {
      // recebemos uma notificação de um cliente ligando p/ a store atual
      handle_incoming_call(data.incoming_call_from_client);
      return;
    }
}

async function handle_incoming_call(call_obj) {
  // ... stuff ...
}
```

```ts
// call_obj:
interface IncomingCallFromClient {
  // ID da ligação
  call_id: string;
  // ID do telefone da store
  phone_id: number;
  // ID da branch da store
  branch_id: string;
  // ID do contacto WABA
  waba_contact_id: string;
  // Número de telefone do contato
  contact_phone_number: string;
  // ID do contato
  contact_id: number;
  // Nome do contato
  contact_name: string;
}
````

### Atendendo uma ligação

Importante: é essencial que esteja atendendo a somente UMA ligação por vez.

```js
const call_storage = {
  pc: null,
  call_id: "",
  contact_id: 0,
  contact_name: "",
  phone_id: 0,
  branch_id: "",
}

async function setup_call(data) {
  // atender a uma ligação envolve duas partes:
  // 1. preparar um WebRTCPeerConnection
  // 2. enviar pro websocket informações sobre a ligação que deseja atender (incluindo o SDP gerado pelo 1)
  
  const ice_config = {
      iceServers: [
          { urls: "stun:stun.l.google.com:19302" },
          { urls: "stun:stun1.l.google.com:3478" },
      ],
  };
  
  call_storage.call_id = data.call_id;
  call_storage.contact_id = data.contact_id;
  call_storage.contact_name = data.contact_name;
  call_storage.phone_id = data.phone_id;
  
  call_storage.pc = new RTCPeerConnection(ice_config);
  
  call_storage.pc.ontrack = (e) => {
    const audio = new Audio();
    audio.srcObject = e.streams[0];
    audio.autoplay = true;
    document.body.appendChild(audio);
    // esta linha acima fará com que o browser pergunte se pode usar o microfone
    
    // obs: no meu pc esse código de exemplo não funciona bem no firefox, mas funciona bem no chrome
  };
  
  // Enviar ICE candidates p/ o backend
  // para negociar a conexão
  call_storage.pc.onicecandidate = (event) => {
      if (event.candidate) {
          console.log("onicecandidate -> ICE candidate", event.candidate);
          if (!event.candidate || !event.candidate.candidate) {
              return;
          }
          ws.send(
              JSON.stringify({
                  type: MSG_SEND_BROWSER_CANDIDATE,
                  send_browser_candidate: {
                      candidate: event.candidate,
                      phone_id: call.phone_id,
                      call_id: call.call_id,
                  },
              }),
          );
      }
  };
  
  const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
  stream.getTracks().forEach((track) => call_storage.pc.addTrack(track, stream));

  const offer = await pc.createOffer();
  await call_storage.pc.setLocalDescription(offer);

  // objeto que será enviado p/ o wabaman via Websockets
  const send_obj = {
      call_id: call_storage.call_id,
      phone_id: call_storage.phone_id,
      branch_id: call_storage.branch_id,
      offer_sdp: offer.sdp,
  };

  ws.send(
      JSON.stringify({
          type: MSG_SETUP_CALL_FROM_BROWSER,
          setup_call_from_browser: send_obj,
      }),
  );
}
```
