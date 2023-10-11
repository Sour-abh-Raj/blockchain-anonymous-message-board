const apiUrl = "http://localhost:8080"; // Adjusted to match the API response structure
const messagesDiv = document.getElementById("messages");
const messageForm = document.getElementById("message-form");
const messageInput = document.getElementById("message-input");

// Function to add a message to the message board
function addMessage(message) {
  const messageDiv = document.createElement("div");
  messageDiv.classList.add("message");
  const messageText = document.createElement("p");
  messageText.textContent = message.data; // Adjusted to match the API response structure
  const timestamp = document.createElement("span");
  timestamp.classList.add("timestamp");
  const isoTimestamp = message.time;
  const formattedTimestamp = new Date(isoTimestamp).toLocaleString();
  timestamp.textContent = formattedTimestamp; // Adjusted to match the API response structure
  messageDiv.appendChild(messageText);
  messageDiv.appendChild(timestamp);
  messagesDiv.appendChild(messageDiv);
}

// Function to fetch and display all messages from the blockchain
function fetchAndDisplayMessages() {
  fetch(`${apiUrl}/messages`)
    .then((response) => response.json())
    .then((data) => {
      data.forEach((message) => {
        addMessage(message);
      });
    })
    .catch((error) => console.error(error));
}

// Initial fetch and display of messages
fetchAndDisplayMessages();

// Function to send a message to the server
function sendMessage(message) {
  // Make a POST request to the server-side endpoint to add a new block to the blockchain
  fetch(`${apiUrl}/addBlock`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      data: message.text,
    }),
  })
    .then((response) => response.json())
    .then((data) => {
      // Once the message is added to the blockchain, add it to the message board
      addMessage(data); // Adjusted to match the API response structure
    })
    .catch((error) => console.error(error));
}

// Event listener for the message form submission
messageForm.addEventListener("submit", (event) => {
  event.preventDefault();
  const message = {
    text: messageInput.value,
    timestamp: Date.now(),
  };
  sendMessage(message);
  messageInput.value = "";
});
