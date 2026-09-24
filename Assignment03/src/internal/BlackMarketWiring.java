package internal;

public class BlackMarketWiring {
    public int shockSystem(float voltage, byte[] data) {
        System.out.println("0x00_INITIATING_SEQUENCE_0xFF");
        System.out.println("0x00_APPLYING_" + voltage + "_VOLTS_0_STREAM_" + data.length + "_BYTES_0xFF");

        if (data.length == 0) {
            return -1;
        }

        if (voltage > 20.0f) {
            return 137;
        }

        System.out.println("0x00_SUCCESS_0xFF");
        return 0;
    }
}
