package internal;

public class HardwareDetector {
    public static ImplantFirmware scanAndMount(String hardwareSignature) {
        System.out.println("[SYSTEM] Scanning Hardware Signature... [SYSTEM]");

        switch(hardwareSignature.toUpperCase()) {
            case "MILITECH_V9":
                System.out.println("[SYSTEM] Verified MilitechOS [SYSTEM]");
                return new MilitechOS();
            case "CORP_ARASAKA_JP":
                System.out.println("[SYSTEM] Verified ArasakaKernel [SYSTEM]");
                return new ArasakaKernel();
            case "ERR_UNKNOWN_UREGISTERED":
                System.out.println("[WARNING] Unregistered Hardware Detected. Mounting Adapter... [WARNING]");
                BlackMarketWiring hardware = new BlackMarketWiring();
                return new BlackMarketAdapter(hardware);
            default:
                throw new IllegalArgumentException("[ERROR] Unrecognized Hardware. Cannot Mount. Aborting... [ERROR]");
        }
    }
}
