import model.NetworkPacketDirector;
import model.NetworkPacketHexDumpBuilder;
import model.NetworkPacketObjectBuilder;

void main() throws UnknownHostException {
    NetworkPacketDirector director = new NetworkPacketDirector();

    NetworkPacketObjectBuilder objBuilder = new NetworkPacketObjectBuilder();
    director.makeHttpPacket(objBuilder, InetAddress.getLocalHost(), InetAddress.getByName("www.google.com"));

    NetworkPacketHexDumpBuilder hexBuilder = new NetworkPacketHexDumpBuilder();
    director.makeDnsPacket(hexBuilder, InetAddress.getLocalHost(), InetAddress.getByName("www.google.com"));

    System.out.println(hexBuilder.getResult());
}
