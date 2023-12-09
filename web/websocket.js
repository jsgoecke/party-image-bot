let socket = new WebSocket("ws://goecke.ngrok.dev/api/v1/ws");

console.log("Attempting websocket connection...")
socket.onopen = function() {
    console.log("Websocket connection established!");
}

socket.onclose = function() {
    console.log("Websocket connection closed!");
}

socket.onmessage = function(message) {
    console.log(message);
    var json = JSON.parse(message.data);

    switch(json.status) {
        case "CONNECTED":
          console.log("Received connect message from websocket");
          break;
        case "IMAGES-GENERATED":
            console.log(json.human_image)
            document.getElementById("statusUpdate").innerHTML = "\"" + json.human_prompt + "\" from " + json.from;
            document.getElementById("aiImage").src = json.ai_image;
            document.getElementById("aiPrompt").innerHTML = json.ai_prompt;
            break;
      }
}