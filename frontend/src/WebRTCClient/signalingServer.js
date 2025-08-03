// const WebSocket = require('ws');
// const http = require('http');

// class SignalingServer {
//   constructor(port = 8080) {
//     this.port = port;
//     this.rooms = new Map(); // roomId -> Set of WebSocket connections
//     this.server = null;
//     this.wss = null;
//   }

//   start() {
//     // Create HTTP server
//     this.server = http.createServer((req, res) => {
//       res.writeHead(200, { 'Content-Type': 'text/plain' });
//       res.end('WebRTC Signaling Server');
//     });

//     // Create WebSocket server
//     this.wss = new WebSocket.Server({ server: this.server });

//     // Handle WebSocket connections
//     this.wss.on('connection', (ws) => {
//       console.log('New client connected');
      
//       let currentRoom = null;

//       ws.on('message', (data) => {
//         try {
//           const message = JSON.parse(data);
//           console.log('Received message:', message.type, 'for room:', message.roomId);

//           switch (message.type) {
//             case 'join-room':
//               this.handleJoinRoom(ws, message.roomId);
//               currentRoom = message.roomId;
//               break;

//             case 'offer':
//             case 'answer':
//             case 'ice-candidate':
//               this.broadcastToRoom(message.roomId, message, ws);
//               break;

//             default:
//               console.log('Unknown message type:', message.type);
//           }
//         } catch (error) {
//           console.error('Error parsing message:', error);
//         }
//       });

//       ws.on('close', () => {
//         console.log('Client disconnected');
//         if (currentRoom) {
//           this.handleLeaveRoom(ws, currentRoom);
//         }
//       });

//       ws.on('error', (error) => {
//         console.error('WebSocket error:', error);
//         if (currentRoom) {
//           this.handleLeaveRoom(ws, currentRoom);
//         }
//       });
//     });

//     // Start server
//     this.server.listen(this.port, () => {
//       console.log(`Signaling server running on port ${this.port}`);
//       console.log(`WebSocket endpoint: ws://localhost:${this.port}`);
//     });
//   }

//   handleJoinRoom(ws, roomId) {
//     // Create room if it doesn't exist
//     if (!this.rooms.has(roomId)) {
//       this.rooms.set(roomId, new Set());
//     }

//     const room = this.rooms.get(roomId);
//     room.add(ws);

//     console.log(`Client joined room: ${roomId}. Total clients in room: ${room.size}`);

//     // Notify other clients in the room
//     this.broadcastToRoom(roomId, {
//       type: 'user-joined',
//       roomId: roomId
//     }, ws);

//     // Send confirmation to the joining client
//     ws.send(JSON.stringify({
//       type: 'room-joined',
//       roomId: roomId,
//       clientCount: room.size
//     }));
//   }

//   handleLeaveRoom(ws, roomId) {
//     const room = this.rooms.get(roomId);
//     if (room) {
//       room.delete(ws);
      
//       // Remove room if empty
//       if (room.size === 0) {
//         this.rooms.delete(roomId);
//         console.log(`Room ${roomId} deleted (empty)`);
//       } else {
//         console.log(`Client left room: ${roomId}. Remaining clients: ${room.size}`);
//       }
//     }
//   }

//   broadcastToRoom(roomId, message, excludeWs = null) {
//     const room = this.rooms.get(roomId);
//     if (room) {
//       room.forEach((client) => {
//         if (client !== excludeWs && client.readyState === WebSocket.OPEN) {
//           try {
//             client.send(JSON.stringify(message));
//           } catch (error) {
//             console.error('Error sending message to client:', error);
//           }
//         }
//       });
//     }
//   }

//   stop() {
//     if (this.wss) {
//       this.wss.close();
//     }
//     if (this.server) {
//       this.server.close();
//     }
//     console.log('Signaling server stopped');
//   }

//   getStats() {
//     const stats = {
//       totalRooms: this.rooms.size,
//       rooms: {}
//     };

//     this.rooms.forEach((clients, roomId) => {
//       stats.rooms[roomId] = clients.size;
//     });

//     return stats;
//   }
// }

// // Export for use as module
// if (typeof module !== 'undefined' && module.exports) {
//   module.exports = SignalingServer;
// }

// // Auto-start if run directly
// if (require.main === module) {
//   const server = new SignalingServer();
//   server.start();

//   // Graceful shutdown
//   process.on('SIGINT', () => {
//     console.log('\nShutting down signaling server...');
//     server.stop();
//     process.exit(0);
//   });

//   // Log stats every 30 seconds
//   setInterval(() => {
//     const stats = server.getStats();
//     console.log('Server stats:', stats);
//   }, 30000);
// } 