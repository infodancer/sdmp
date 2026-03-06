package sdmp

import (
	"fmt"
	"net"

	pb "github.com/infodancer/sdmp/proto/sdmp/v1"
	"google.golang.org/protobuf/proto"
)

// NotificationListener handles incoming UDP notification packets.
// Notifications are fire-and-forget — no response is ever sent.
type NotificationListener struct {
	conn *net.UDPConn
}

// NewNotificationListener creates a UDP listener on the given address.
func NewNotificationListener(addr string) (*NotificationListener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("resolve udp address: %w", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("listen udp: %w", err)
	}
	return &NotificationListener{conn: conn}, nil
}

// Listen reads notification packets from the UDP socket. It calls the
// handler function for each successfully parsed notification. This method
// blocks until the connection is closed.
func (l *NotificationListener) Listen(handler func(*pb.EncryptedNotification)) error {
	buf := make([]byte, 65535)
	for {
		n, _, err := l.conn.ReadFromUDP(buf)
		if err != nil {
			return fmt.Errorf("read udp: %w", err)
		}

		var enc pb.EncryptedNotification
		if err := proto.Unmarshal(buf[:n], &enc); err != nil {
			// Malformed packets are silently dropped — no response.
			continue
		}

		handler(&enc)
	}
}

// Close shuts down the UDP listener.
func (l *NotificationListener) Close() error {
	return l.conn.Close()
}

// NotificationSender sends UDP notification packets to receiving domains.
type NotificationSender struct{}

// Send serializes and sends an encrypted notification to the given address.
func (s *NotificationSender) Send(addr string, notification *pb.EncryptedNotification) error {
	data, err := proto.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("resolve udp address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("dial udp: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write udp: %w", err)
	}
	return nil
}
