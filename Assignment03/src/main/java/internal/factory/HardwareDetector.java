package internal.factory;

import external.BlackMarketWiring;
import internal.adapter.BlackMarketAdapter;
import internal.implementor.ArasakaKernel;
import internal.implementor.ImplantFirmware;
import internal.implementor.MilitechOS;

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
            case "ERR_UNKNOWN_UNREGISTERED":
                System.out.println("[WARNING] Unregistered Hardware Detected. Mounting Adapter... [WARNING]");
                BlackMarketWiring hardware = new BlackMarketWiring();
                return new BlackMarketAdapter(hardware);
            default:
                System.out.println("[ERROR] Unrecognized Hardware. Cannot Mount. Aborting... [ERROR]");
                System.exit(0);
                return null;
        }
    }
}
