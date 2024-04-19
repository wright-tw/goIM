$(document).ready(function () {
  const ACTION_SYSTEM_MSG = 1;
  const ACTION_MSG = 2;
  const ACTION_CHANGE_NAME = 3;
  const ACTION_ONLINE_PEOPLE = 4;

  let username = getUsername();
  let currentDomain = window.location.hostname;
  let currentPort = window.location.port;
  let currentProtocol = window.location.protocol;
  let wsProtocol = "ws";
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

  // 處理接收到的訊息
  socket.onmessage = function (event) {
    handleMessage(event);
  };

  // 當 WebSocket 連接關閉時，清除心跳定時器
  socket.onclose = function (event) {
    clearInterval(heartbeatTimer);
  };

  // 當 WebSocket 連接打開時的處理函數
  socket.onopen = function (event) {
    console.log("WebSocket connected");
  };

  // 處理接收到的訊息
  function handleMessage(event) {
    if (event.data == "pong") {
      return;
    }
    const data = JSON.parse(event.data);
    const action = data.action;
    switch (action) {
      case ACTION_MSG:
      case ACTION_SYSTEM_MSG:
        // 如果是普通訊息，則處理訊息內容
        handleRegularMessage(data);
        break;
      case ACTION_ONLINE_PEOPLE:
        $("#online_people_count").html(data.msg);
        break;
      default:
        console.log("Unknown action:", action);
    }
  }

  // 處理普通訊息
  function handleRegularMessage(data) {
    const msgUsername = data.username ?? "";
    const msgString = data.msg ?? "";
    const timeString = data.time ?? "";
    const messageContainer = createMessageContainer(
      msgUsername,
      msgString,
      timeString
    );
    $("#chat-messages").append(messageContainer);
    $("#chat-messages").scrollTop($("#chat-messages")[0].scrollHeight);
  }

  function createMessageContainer(username, message, timeString) {
    const avatar = generateAvatarFromString(username);
    const messageDiv = $("<div>")
      .addClass("message-content")
      .text(username + " : " + message);
    const timeDiv = $("<div>")
      .addClass("time-content")
      .text(timeString);
    const messageContainer = $("<div>")
      .addClass("message-container")
      .append(avatar, messageDiv, timeDiv);
    return messageContainer;
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

  // 獲取用戶名
  function getUsername() {
    let storedUsername = localStorage.getItem("username") ?? "";
    let username = "";
    if (storedUsername.trim() == "") {
      username = prompt("請輸入聊天室想用的名稱:");
      if (username.trim() !== "") {
        localStorage.setItem("username", username);
      }
    } else {
      username = storedUsername;
    }
    return username;
  }

  // 根據字串生成顏色碼
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

  // 心跳函數
  function heartbeat() {
    socket.send("ping");
  }

  // 發送消息函數
  function sendMessage() {
    const messageInput = $("#message-input");
    const message = messageInput.val();
    if (message.trim() !== "") {
      socket.send(message);
      messageInput.val("");
    }
  }

  // 監聽發送按鈕點擊事件
  $("#send-button").on("click", sendMessage);

  // 監聽輸入框按下 Enter 鍵事件
  $("#message-input").on("keypress", function (event) {
    if (event.keyCode === 13) {
      sendMessage();
    }
  });

  // 監聽更換用戶名按鈕點擊事件
  $("#change-username-button").on("click", function () {
    // 提示用戶輸入新的用戶名
    let newUsername = prompt("請重新輸入想要的聊天室用戶名:");
    if (newUsername.trim() !== "") {
      // 更新 localStorage 中的用戶名
      localStorage.setItem("username", newUsername);
      // 重新連接 WebSocket 並使用新的用戶名
      reconnectWebSocket(newUsername);
    }
  });

  // 重新連接 WebSocket 的方法
  function reconnectWebSocket(newUsername) {
    // 關閉原來的 WebSocket 連接
    socket.close();

    // 建立新的 WebSocket 連接
    socket = new WebSocket(
      wsProtocol +
        "://" +
        currentDomain +
        ":" +
        currentPort +
        "/ws?name=" +
        newUsername +
        "&rename=" +
        1
    );

    // 重新設置 WebSocket 事件處理函數
    socket.onmessage = function (event) {
      handleMessage(event);
    };
    socket.onclose = function (event) {
      clearInterval(heartbeatTimer);
    };
    socket.onopen = function (event) {
      console.log("WebSocket connected");
    };
  }

});
