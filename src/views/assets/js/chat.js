// 提示用戶輸入用戶名
const username = prompt("請輸入聊天室想用的名稱:");

// 如果用戶名不是空的，則將其存儲在localStorage中
if (username.trim() !== "") {
  localStorage.setItem("username", username);
}

// 加載用戶名
// const storedUsername = localStorage.getItem("username");

const currentDomain = window.location.hostname;
const currentPort = window.location.port;
const socket = new WebSocket(
  "wss://" + currentDomain + ":" + currentPort + "/ws?name=" + username
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

  // 使用 Set 來存儲字串中的唯一 Unicode 編碼值
  const uniqueChars = new Set([...str]);

  // 將 Set 轉換為數組，並根據 Unicode 編碼值排序
  const sortedUniqueChars = Array.from(uniqueChars).sort();

  // 迴圈遍歷排序後的唯一字元
  for (const char of sortedUniqueChars) {
    // 取得字元的 Unicode 編碼值
    const charCode = char.charCodeAt(0);
    // 根據字元的編碼值計算對應的色碼部分
    const colorPart = (charCode % 255).toString(16); // 將編碼值取模 255，再轉換為十六進制字串
    // 將色碼部分添加到色碼中
    colorCode += colorPart.padStart(2, "0"); // 補零至兩位
  }

  // 返回生成的色碼
  return colorCode;
}

