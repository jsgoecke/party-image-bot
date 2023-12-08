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
        case "SMS-RECEIVED":
          var msg = "SMS Recived from " + json.from + " with prompt: " + "\"" + json.human_prompt + "\"";
          document.getElementById("statusUpdate").innerHTML = msg;
          document.getElementById("humanImage").src = "https://goecke.ngrok.dev/web/cat-spin.gif";
          document.getElementById("aiImage").src = "https://goecke.ngrok.dev/web/cat-spin.gif";
          break;
        case "IMAGES-GENERATED":
            console.log(json.human_image)
            document.getElementById("statusUpdate").innerHTML = "Images by: " + json.from;
            document.getElementById("humanImage").src = json.human_image;
            document.getElementById("aiImage").src = json.ai_image;
            document.getElementById("humanPrompt").innerHTML = json.human_prompt;
            document.getElementById("aiPrompt").innerHTML = json.ai_prompt;
            break;
      }
}