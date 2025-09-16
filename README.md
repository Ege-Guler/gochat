# GoChat - Peer-to-Peer Encrypted Chat Application

GoChat is a secure, peer-to-peer chat application written in Go that combines multiple networking technologies to enable direct communication between users while bypassing network restrictions.

## What is this about?

This project implements a sophisticated P2P chat system that uses:

- **NAT Traversal**: Uses STUN (Session Traversal Utilities for NAT) to discover public IP addresses and enable direct communication through NAT devices
- **UDP Hole Punching**: Establishes direct peer-to-peer connections even when both peers are behind NAT
- **MQTT Relay**: Uses MQTT broker as a fallback relay system for initial connection establishment and key exchange
- **End-to-End Encryption**: Implements ECDH (Elliptic Curve Diffie-Hellman) key exchange with AES-GCM encryption for secure messaging
- **Real-time Chat**: Provides a terminal-based chat interface for direct communication

## Architecture

### Core Components

1. **Peer-to-Peer Chat (`internal/chat/peer.go`)**
   - Direct UDP communication between peers
   - Message queuing and routing
   - Terminal-based chat interface

2. **Network Configuration (`netconf/netconfig.go`)**
   - STUN client for NAT traversal
   - Public IP discovery using Google's STUN server
   - UDP connection management

3. **MQTT Relay System (`internal/mqtt/`)**
   - Session management and peer discovery
   - Secure key exchange messaging
   - Fallback communication channel

4. **Cryptographic Security (`internal/crypto/crypto.go`)**
   - ECDH key exchange for forward secrecy
   - AES-256-GCM encryption for message protection
   - Secure random key generation

5. **Configuration Management (`internal/config/`)**
   - YAML-based configuration loading
   - MQTT broker settings
   - User preferences

### Communication Flow

1. **Initial Setup**: Each peer generates ECDH key pairs and creates a unique session ID
2. **NAT Discovery**: Uses STUN to discover public IP address and port
3. **Peer Discovery**: Exchanges public keys via MQTT broker
4. **Connection Establishment**: Performs UDP hole punching to establish direct connection
5. **Secure Communication**: All messages are encrypted with session-specific AES keys

## Features

- ✅ **NAT Traversal**: Works behind most NAT configurations
- ✅ **End-to-End Encryption**: Messages are encrypted using ECDH + AES-GCM
- ✅ **Direct P2P**: No central server required for messaging (only for initial discovery)
- ✅ **Real-time**: Low-latency UDP-based communication
- ✅ **Terminal Interface**: Simple command-line chat interface
- ✅ **Session Management**: Unique session IDs prevent message interception

## Technical Details

### Dependencies

- **MQTT Client**: `github.com/eclipse/paho.mqtt.golang` for broker communication
- **STUN Library**: `github.com/pion/stun` for NAT traversal
- **UUID Generation**: `github.com/google/uuid` for session management
- **YAML Parsing**: `gopkg.in/yaml.v3` for configuration

### Security Features

- **Forward Secrecy**: New ECDH keys for each session
- **Authenticated Encryption**: AES-GCM provides both confidentiality and integrity
- **Session Isolation**: Each chat session has unique cryptographic keys
- **Secure Key Exchange**: ECDH prevents man-in-the-middle attacks on key exchange

## Usage

### Prerequisites

1. Go 1.22.2 or later
2. Access to an MQTT broker
3. Network connectivity (UDP ports)

### Configuration

Create a `conf.yaml` file with your MQTT broker details:

```yaml
mqtt:
  broker: "your-mqtt-broker.com"
  port: 8883
  topic: "gochat/exchange"
  clientID: "your-client-id"
  username: "your-username"
  password: "your-password"
  qos: 1
  retain: false

chat:
  username: "your-chat-username"
```

### Running the Application

```bash
# Build the application
go build -o gochat ./cmd/peer

# Run the chat client
./gochat
```

The application will:
1. Display your public and local IP addresses
2. Prompt for the remote peer's address (IP:PORT)
3. Attempt to establish a direct P2P connection
4. Start the chat interface once connected

## Development Status

This appears to be a work-in-progress project with some features still under development:

- ✅ Basic P2P chat functionality
- ✅ NAT traversal and hole punching
- ✅ ECDH key exchange
- ✅ AES encryption/decryption
- 🚧 Complete MQTT relay implementation
- 🚧 Error handling and reconnection logic
- 🚧 Configuration validation
- 🚧 Comprehensive testing

## Network Requirements

- UDP port access for peer-to-peer communication
- Access to STUN servers (uses `stun.l.google.com:19302`)
- Access to configured MQTT broker for initial handshake
- Firewall rules allowing UDP traffic on dynamic ports

This project demonstrates advanced networking concepts including NAT traversal, cryptographic protocols, and distributed system design in a practical peer-to-peer chat application.