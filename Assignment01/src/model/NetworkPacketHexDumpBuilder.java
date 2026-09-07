package model;

import java.net.InetAddress;
import java.util.Objects;

public class NetworkPacketHexDumpBuilder implements PacketBuilder {

    private InetAddress sourceIP;
    private InetAddress destinationIP;
    private Integer sourcePort;
    private Integer destinationPort;
    private Protocol protocol;
    private String payload;

    public NetworkPacketHexDumpBuilder setSourceIP(InetAddress sourceIP) {
        Objects.requireNonNull(sourceIP, "Source IP cannot be null");
        this.sourceIP = sourceIP;
        return this;
    }

    public NetworkPacketHexDumpBuilder setDestinationIP(InetAddress destinationIP) {
        Objects.requireNonNull(destinationIP, "Destination IP cannot be null");
        this.destinationIP = destinationIP;
        return this;
    }

    public NetworkPacketHexDumpBuilder setSourcePort(Integer sourcePort) {
        Objects.requireNonNull(sourcePort, "Source port cannot be null");
        checkPortBounds(sourcePort);
        this.sourcePort = sourcePort;
        return this;
    }

    public NetworkPacketHexDumpBuilder setDestinationPort(Integer destinationPort) {
        Objects.requireNonNull(destinationPort, "Destination port cannot be null");
        checkPortBounds(destinationPort);
        this.destinationPort = destinationPort;
        return this;
    }

    public NetworkPacketHexDumpBuilder setProtocol(Protocol protocol) {
        Objects.requireNonNull(protocol, "Protocol cannot be null");
        this.protocol = protocol;
        return this;
    }

    public NetworkPacketHexDumpBuilder setPayload(String payload) {
        Objects.requireNonNull(payload, "Payload cannot be null");
        this.payload = payload;
        return this;
    }

    public String getResult() {
        if (Objects.isNull(sourceIP) || Objects.isNull(destinationIP) || Objects.isNull(protocol)) {
            throw new IllegalStateException("Source IP, Destination IP, and Protocol must be set");
        }


        sourcePort = (sourcePort != null) ? sourcePort : 80;
        destinationPort = (destinationPort != null) ? destinationPort : 80;
        payload = (payload != null) ? payload : "";

        final StringBuilder sb = new StringBuilder();
        sb.append("=== NETWORK PACKET SPEC ===\n");
        sb.append("Source IP: ").append(sourceIP.getHostAddress()).append("\n");
        sb.append("Destination IP: ").append(destinationIP.getHostAddress()).append("\n");
        sb.append("Source Port: ").append(sourcePort).append("\n");
        sb.append("Destination Port: ").append(destinationPort).append("\n");
        sb.append("Protocol: ").append(protocol).append("\n");
        sb.append("\n=== HEX DUMP ===\n");

        for (byte b : sourceIP.getAddress()) {
            sb.append(String.format("%02X ", b));
        }

        for (byte b : destinationIP.getAddress()) {
            sb.append(String.format("%02X ", b));
        }

        for (byte b : payload.getBytes()) {
            sb.append(String.format("%02X ", b));
        }

        return sb.append("\n").toString();
    }
}
