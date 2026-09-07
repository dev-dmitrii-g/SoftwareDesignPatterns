package model;

import java.net.InetAddress;

public interface PacketBuilder {
    PacketBuilder setSourceIP(InetAddress sourceIP);
    PacketBuilder setDestinationIP(InetAddress destinationIP);
    PacketBuilder setSourcePort(Integer sourcePort);
    PacketBuilder setDestinationPort(Integer destinationPort);
    PacketBuilder setProtocol(Protocol protocol);
    PacketBuilder setPayload(String payload);

    default void checkPortBounds(Integer port) {
        if (port < 0 || port > 65535) {
            throw new IllegalArgumentException("Port must be between 0 and 65535");
        }
    }
}