package client;

import internal.asbtraction.MotorOverride;
import internal.asbtraction.NeuralInterface;
import internal.exception.RejectionSyndromeException;
import internal.factory.HardwareDetector;
import internal.implementor.ImplantFirmware;

import java.util.Scanner;

public class Boot {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        System.out.println("Choose the manufacturer of the implant: ");
        System.out.println("1. Militech");
        System.out.println("2. Arasaka Corporation");
        System.out.println("3. Unregistered");
        System.out.println("4. Other");
        int hardwareSignature = scanner.nextInt();

        String hardwarePort = switch (hardwareSignature) {
            case 1 -> "MILITECH_V9";
            case 2 -> "CORP_ARASAKA_JP";
            case 3 -> "ERR_UNKNOWN_UNREGISTERED";
            default -> "UNKNOWN";
        };

        try {
            ImplantFirmware firmware = HardwareDetector.scanAndMount(hardwarePort);
            NeuralInterface override = new MotorOverride(firmware);
            override.triggerImplant();
        } catch (RejectionSyndromeException e) {
            System.out.println("[ALERT] " + e.getMessage() + " [ALERT]");
        }
    }
}
