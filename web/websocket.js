let socket = new WebSocket("ws://goecke.ngrok.dev/api/v1/ws");

console.log("Attempting websocket connection...")
socket.onopen = function() {
    console.log("Websocket connection established!");
}

socket.onclose = () => {
    console.log('WebSocket connection closed');
    if (retries > 0) {
        retries--;
        console.log('Retrying websocket connection...)');
        setTimeout(connectWebSocket, 5000); // Retry after 5 seconds
    }
};

socket.onmessage = function(message) {
    console.log(message);
    var json = JSON.parse(message.data);

    switch(json.status) {
        case "CONNECTED":
          console.log("Received connect message from websocket");
          break;
        case "IMAGES-GENERATED":
            if (json.ai_image != null){
                document.getElementById("aiImage").src = json.ai_image;
            }
            break;
      }
}