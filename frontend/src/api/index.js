var socket = new WebSocket("ws://127.0.0.1:8081/ws");

let connect = (callback) => {
  console.log("Attempting Connection...");

  socket.onopen = (event) => {
    console.log("socket.onopen: Successfully Connected:", event);
  };

  socket.onmessage = msg => {
    console.log("socket.onmessage: ", msg);
    callback(msg)
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("sending msg: ", msg);
  socket.send(msg);
};

export { connect, sendMsg };