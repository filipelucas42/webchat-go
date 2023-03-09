let socket
let board = document.getElementById("board")
let username = document.getElementById("username")
let createUser = document.getElementById("createUser")
let closeConnection = document.getElementById("closeConnection")
let sendMessage = document.getElementById("chat")

function SocketOpen(){
    createUser.style.display = "none"
    closeConnection.style.display = "block"
    sendMessage.style.display = "block"
    username.innerHTML = "Your username: " + document.getElementById("name").value
}

function SocketClose(){
    createUser.style.display = "block"
    closeConnection.style.display = "none"
    sendMessage.style.display = "none"
    document.getElementById("name").value=""
    username.innerHTML = ""
}
closeConnection.addEventListener("click", (e)=>{
    socket.close()
    createUser.style.display = "none"
    closeConnection.style.display = "none"
})
document.getElementById("chat").addEventListener("submit", (event)=>{
    event.preventDefault()
    text = document.getElementById("text")
    socket.send(text.value)
    text.value=""
})
createUser.addEventListener('submit', (e)=>{
    e.preventDefault()
    socket = new WebSocket("ws://127.0.0.1:8082/socket")

    socket.onopen = function(e) {
        socket.send(document.getElementById("name").value)
        SocketOpen()
    };

    socket.onmessage = function(event) {
        board.value = board.value + "\n" + event.data
        board.scrollTop = board.scrollHeight
    };

    socket.onclose = function(event) {
        board.value = board.value + "\n" + "Connection to chat closed"
        SocketClose()
    };

    socket.onerror = function(error) {
    };

})