const express = require('express');
const http = require('http');
const app = express();
const server = http.createServer(app);
const io = require('socket.io')(server);

// Servir les fichiers statiques depuis le répertoire public
app.use(express.static('public'));

// Logique de gestion des connexions WebSocket
io.on('connection', (socket) => {
    console.log('Nouvelle connexion WebSocket');
    
    // Logique de gestion des messages
    socket.on('message', (message) => {
        console.log('Message reçu:', message);
        io.emit('message', message); // Diffuser le message à tous les clients connectés
    });
});

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
    console.log(`Serveur écoutant sur le port ${PORT}`);
});
