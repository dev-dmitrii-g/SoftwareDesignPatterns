package model;

import java.net.InetAddress;

public class NetworkPacketDirector {
    public void makeHttpPacket(PacketBuilder builder, InetAddress src, InetAddress dst) {
        builder.setSourceIP(src)
                .setDestinationIP(dst)
                .setSourcePort(80)
                .setDestinationPort(80)
                .setProtocol(Protocol.TCP)
                .setPayload("GET / HTTP/1.1");
    }

    public void makeDnsPacket(PacketBuilder builder, InetAddress src, InetAddress dst) {
        builder.setSourceIP(src)
                .setDestinationIP(dst)
                .setSourcePort(53)
                .setDestinationPort(53)
                .setProtocol(Protocol.UDP)
                .setPayload("example.com");
    }
}
