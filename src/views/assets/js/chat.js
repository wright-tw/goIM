let username = getUsername();
let currentDomain = window.location.hostname;
let currentPort = window.location.port;
let currentProtocol = window.location.protocol;
let wsProtocol = "ws";
console.log(currentProtocol);
if (currentProtocol == "https:") {
  wsProtocol = "wss";
}
let socket = new WebSocket(
  wsProtocol +
    "://" +
    currentDomain +
    ":" +
    currentPort +
    "/ws?name=" +
    username
);

// 啟動心跳定時器
const heartbeatTimer = setInterval(heartbeat, 5000);

// 建立心跳方法
function heartbeat() {
  // 發送心跳訊號到服務器（這裡假設使用 socket 變數代表 WebSocket 連接）
  socket.send("ping");
}

// 當 WebSocket 連接關閉時，清除心跳定時器
socket.onclose = function (event) {
  clearInterval(heartbeatTimer);
};

socket.onopen = function (event) {
  console.log("WebSocket connected");
};

socket.onmessage = function (event) {
  let msg = event.data;
  let parts = msg.split(":");
  let msgUsername = parts[0] ?? "";
  let msgString = parts[1] ?? "";

  if (msg == "pong") {
    return;
  }

  // 創建一個包含頭像和訊息的 div 元素
  const messageContainer = document.createElement("div");
  messageContainer.classList.add("message-container");

  // 製作頭像
  let avatar = generateAvatarFromString(msgUsername);
  messageContainer.appendChild(avatar);

  // 製作訊息
  const messageDiv = document.createElement("div");
  messageDiv.textContent = msg;
  messageDiv.classList.add("message-content"); // 添加 message-content 類
  messageContainer.appendChild(messageDiv);

  // 設置頭像和訊息之間的距離
  messageContainer.style.marginBottom = "10px";

  // 添加整個訊息容器到聊天視窗中
  document.getElementById("chat-messages").appendChild(messageContainer);

  // 滾動到底部
  document.getElementById("chat-messages").scrollTop =
    document.getElementById("chat-messages").scrollHeight;
};

// 發送消息函數
function sendMessage() {
  const messageInput = document.getElementById("message-input");
  const message = messageInput.value;
  if (message.trim() !== "") {
    socket.send(message);
    messageInput.value = "";
  }
}

// 監聽按鍵事件
document.getElementById("send-button").addEventListener("click", sendMessage);

document
  .getElementById("message-input")
  .addEventListener("keypress", function (event) {
    // 如果按下的是 Enter 鍵，則發送消息
    if (event.keyCode === 13) {
      sendMessage();
    }
  });

function getUsername() {
  // 加載用戶名
  let storedUsername = localStorage.getItem("username") ?? "";
  let username = "";
  if (storedUsername.trim() == "") {
    // 提示用戶輸入用戶名
    username = prompt("請輸入聊天室想用的名稱:");

    // 如果用戶名不是空的，則將其存儲在localStorage中
    if (username.trim() !== "") {
      localStorage.setItem("username", username);
    }
  } else {
    username = storedUsername;
  }

  return username;
}

function generateAvatarFromString(text) {
  // 創建一個用於繪製頭像的 Canvas 元素
  const canvas = document.createElement("canvas");
  const size = 45; // 設置 Canvas 寬度和高度
  canvas.width = size;
  canvas.height = size;

  // 獲取 2D 繪圖上下文
  const ctx = canvas.getContext("2d");

  // 根據用戶名計算顏色
  const color = getColorCodeByString(text);

  // 設置填充顏色為用戶名計算出的顏色
  ctx.fillStyle = color;

  // 繪製正方形頭像
  ctx.fillRect(0, 0, size, size); // 繪製矩形作為頭像背景

  // 在頭像上繪製文字（只顯示第一個文字）
  ctx.fillStyle = "#ffffff"; // 文字顏色為白色
  ctx.font = "30px Arial"; // 設置字體大小和字型
  ctx.textAlign = "center"; // 文字居中對齊
  ctx.textBaseline = "middle"; // 文字垂直居中對齊
  ctx.fillText(text[0], size / 2, size / 2); // 只繪製第一個文字，位置居中

  // 返回生成的 Canvas 對象
  return canvas;
}

function hashCode(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  return hash;
}

function intToRGB(i) {
  const c = (i & 0x00ffffff).toString(16).toUpperCase();
  return "#" + "00000".substring(0, 6 - c.length) + c;
}

function getColorCodeByString(str) {
  return intToRGB(hashCode(str));
}
