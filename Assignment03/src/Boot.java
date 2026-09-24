import internal.*;

public class Boot {
    public static void main(String[] args) {
        String hardwarePort = args[0].toUpperCase();

        try {
            ImplantFirmware firmware = HardwareDetector.scanAndMount(hardwarePort);
            NeuralInterface override = new MotorOverride(firmware);
            override.triggerImplant();
        } catch (RejectionSyndromeException e) {
            System.out.println("[ALERT] " + e.getMessage() + " [ALERT]");
        }
    }
}
