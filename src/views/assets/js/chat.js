// 提示用戶輸入用戶名
const username = prompt("Please enter your username:");

// 如果用戶名不是空的，則將其存儲在localStorage中
if (username.trim() !== "") {
  localStorage.setItem("username", username);
}

// 加載用戶名
// const storedUsername = localStorage.getItem("username");

const currentDomain = window.location.hostname;
const currentPort = window.location.port;
const socket = new WebSocket(
  "ws://" + currentDomain + ":" + currentPort + "/ws?name=" + username
);

socket.onopen = function (event) {
  console.log("WebSocket connected");
};

socket.onmessage = function (event) {
  const messageDiv = document.createElement("div");
  messageDiv.textContent = event.data;
  let parts = messageDiv.textContent.split(":");
  let msgUsername = parts[0];

  messageDiv.style.color = getColorCodeByString(msgUsername);
  messageDiv.classList.add("message"); // 添加 message 類
  document.getElementById("chat-messages").appendChild(messageDiv);

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

function getColorCodeByString(str) {
  // 初始化色碼值
  let colorCode = "#";

  // 迴圈遍歷字串的每個字元
  for (let i = 0; i < str.length; i++) {
    // 取得字元的 Unicode 編碼值
    const charCode = str.charCodeAt(i);
    // 將編碼值轉換為十六進制字串
    const hexCode = charCode.toString(16);
    // 將十六進制字串補零至兩位
    const paddedHexCode = hexCode.padStart(2, "0");
    // 將補零後的十六進制字串添加到色碼中
    colorCode += paddedHexCode;
  }

  // 返回生成的色碼
  return colorCode;
}
