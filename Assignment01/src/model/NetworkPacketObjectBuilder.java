package model;

import java.net.InetAddress;
import java.util.Objects;

public class NetworkPacketObjectBuilder implements PacketBuilder {
    private InetAddress sourceIP;
    private InetAddress destinationIP;
    private Integer sourcePort;
    private Integer destinationPort;
    private Protocol protocol;
    private String payload;

    public NetworkPacketObjectBuilder setSourceIP(InetAddress sourceIP) {
        Objects.requireNonNull(sourceIP, "Source IP cannot be null");
        this.sourceIP = sourceIP;
        return this;
    }

    public NetworkPacketObjectBuilder setDestinationIP(InetAddress destinationIP) {
        Objects.requireNonNull(destinationIP, "Destination IP cannot be null");
        this.destinationIP = destinationIP;
        return this;
    }

    public NetworkPacketObjectBuilder setSourcePort(Integer sourcePort) {
        Objects.requireNonNull(sourcePort, "Source port cannot be null");
        checkPortBounds(sourcePort);
        this.sourcePort = sourcePort;
        return this;
    }

    public NetworkPacketObjectBuilder setDestinationPort(Integer destinationPort) {
        Objects.requireNonNull(destinationPort, "Destination port cannot be null");
        checkPortBounds(destinationPort);
        this.destinationPort = destinationPort;
        return this;
    }

    public NetworkPacketObjectBuilder setProtocol(Protocol protocol) {
        Objects.requireNonNull(protocol, "Protocol cannot be null");
        this.protocol = protocol;
        return this;
    }

    public NetworkPacketObjectBuilder setPayload(String payload) {
        Objects.requireNonNull(payload, "Payload cannot be null");
        this.payload = payload;
        return this;
    }

    public NetworkPacket getResult() {
        if (Objects.isNull(sourceIP) || Objects.isNull(destinationIP) || Objects.isNull(protocol)) {
            throw new IllegalStateException("Source IP, Destination IP, and Protocol must be set");
        }

        if (sourcePort == null) {
            sourcePort = 80;
        }

        if (destinationPort == null) {
            destinationPort = 80;
        }

        if (payload == null) {
            payload = "";
        }

        return new NetworkPacket(sourceIP, destinationIP, sourcePort, destinationPort, protocol, payload);
    }

    private void checkPortBounds(Integer port) {
        if (port < 0 || port > 65535) {
            throw new IllegalArgumentException("Port must be between 0 and 65535");
        }
    }
}
