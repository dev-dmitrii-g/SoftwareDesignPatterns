package model;

import java.net.InetAddress;

public interface PacketBuilder {
    PacketBuilder setSourceIP(InetAddress sourceIP);
    PacketBuilder setDestinationIP(InetAddress destinationIP);
    PacketBuilder setSourcePort(Integer sourcePort);
    PacketBuilder setDestinationPort(Integer destinationPort);
    PacketBuilder setProtocol(Protocol protocol);
    PacketBuilder setPayload(String payload);
}