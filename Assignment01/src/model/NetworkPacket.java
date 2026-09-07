package model;

import java.net.InetAddress;

public class NetworkPacket {
    private final InetAddress sourceIP;
    private final InetAddress destinationIP;
    private final Integer sourcePort;
    private final Integer destinationPort;
    private final Protocol protocol;
    private final String payload;

    NetworkPacket(InetAddress sourceIP, InetAddress destinationIP, int sourcePort, int destinationPort, Protocol protocol, String payload) {
        this.sourceIP = sourceIP;
        this.destinationIP = destinationIP;
        this.sourcePort = sourcePort;
        this.destinationPort = destinationPort;
        this.protocol = protocol;
        this.payload = payload;
    }

    public InetAddress getSourceIP() {
        return sourceIP;
    }

    public InetAddress getDestinationIP() {
        return destinationIP;
    }

    public Integer getSourcePort() {
        return sourcePort;
    }

    public Integer getDestinationPort() {
        return destinationPort;
    }

    public Protocol getProtocol() {
        return protocol;
    }

    public String getPayload() {
        return payload;
    }

    @Override
    public String toString() {
        return "NetworkPacket{" + "sourceIP=" + sourceIP + ", destinationIP=" + destinationIP + ", sourcePort=" + sourcePort + ", destinationPort=" + destinationPort + ", protocol=" + protocol + ", payload=" + payload + '}';
    }
}
