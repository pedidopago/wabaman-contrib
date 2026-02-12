# Ligações

## Ligações Iniciadas pelo Cliente

```js
// ws = new WebSocket("path-to-wabaman-wss-endpoint");

// eventos relevantes:
const MSG_INCOMING_CALL_FROM_CLIENT = 30;
const MSG_SETUP_CALL_FROM_BROWSER = 31;
const MSG_TERMINATE_CALL = 32;
const MSG_CALL_CONSUMED = 33;
const MSG_ACCEPT_CALL = 34; //TODO(gabs): revisar este evento
const MSG_REJECT_CALL = 35; //TODO(gabs): revisar este evento
const MSG_SEND_BROWSER_CANDIDATE = 36;
const MSG_CALL_STARTED = 37; //TODO(gabs): revisar este evento
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
  started_at: 0,
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

### Recebendo SDP do servidor

Recebendo um SDP do servidor Wabaman p/ estabelecer a conexão webrtc.

```js
ws.onmessage = function (event) {
  const data = JSON.parse(event.data);
  
  if (data.type === MSG_CALL_ON_ANSWER_SDP) {
    // recebemos os dados necessários para estabelecer a conexão
    handle_establish_call(data.call_on_answer_sdp);
    return;
  }
  
  if (data.type === MSG_CALL_START_TIMER) {
    // este evento é recebido após o MSG_CALL_ON_ANSWER_SDP, quando a conexão foi extabelecida
    // com sucesso nas duas pontas
  }
}

async function handle_establish_call(call_on_answer_sdp) {
  await pc.setRemoteDescription(new RTCSessionDescription({ type: "answer", sdp: call_on_answer_sdp.sdp }));
}

async function on_call_start_timer(jobj) {
  call_storage.started_at = Date.now();

  //TODO: mostrar visualmente feedback da ligação iniciada
  
  // exemplo de um timer que atualiza a cada segundo:
  call_storage.timer_interval = setInterval(() => {
      const elapsed = Date.now() - call.started_at;
      const raw_seconds = Math.floor(elapsed / 1000);
      const minutes = Math.floor(raw_seconds / 60);
      const seconds = raw_seconds % 60;
      const txt = `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
      document.getElementById("timer").textContent = txt;
  }, 1000);
}
```

### Encerrando uma Ligação

```js

async function terminate_call() {
  if (call_storage.pc == null) return;
  if (call_storage.call_id == "") return;
  
  ws.send(
    JSON.stringify({
      type: MSG_TERMINATE_CALL,
      terminate_call: {
        call_id: call.call_id,
      },
    }),
  );
  
  //TODO: representar visualmente que a ligação terminou
}
```

### Recebendo Evento que o cliente encerrou a ligação

```js
ws.onmessage = function (event) {
  const data = JSON.parse(event.data);
  
  if (data.type === MSG_CALL_ENDED) {
    if (call_storage.pc != null && call_storage.call_id == data.call_ended.call_id) {
      // checar se a ligação é a ligação que está ativa pra esse agent no momento:
      on_call_ended(data.call_ended);
    }
    return;
  }
}

async function on_call_ended(call_ended) {
  if (call_storage.pc == null) return;
  pc.close();

  // Remove o component audio adicionado anteriormente
  document.querySelectorAll("audio").forEach((audio) => audio.remove());

  call_storage.branch_id = "";
  call_storage.call_id = "";
  call_storage.contact_id = 0;
  call_storage.contact_name = "";
  call_storage.contact_phone_number = "";
  call_storage.pc = null;
  call_storage.phone_id = 0;
  call_storage.started_at = 0;
  call_storage.waba_contact_id = "";
  try {
      // no exemplo da ligação sendo efetivada, criamos esse timer 
      clearInterval(call_storage.timer_interval);
      call_storage.timer_interval = null;
  } catch (error) {
      console.error(`clearInterval failed ${error}`);
  }

  //TODO: alterar UI que mostra visualmente que a ligação terminou
}
```
