package internal.adapter;

import external.BlackMarketWiring;
import internal.implementor.NeuralPayload;
import internal.exception.RejectionSyndromeException;
import internal.implementor.ImplantFirmware;

import java.nio.charset.StandardCharsets;

public class BlackMarketAdapter implements ImplantFirmware {
    private final BlackMarketWiring hardware;
    public BlackMarketAdapter(BlackMarketWiring hardware) {
        this.hardware = hardware;
    }

    @Override
    public void executeNeuralCommand(NeuralPayload payload) throws RejectionSyndromeException {
        float voltage = payload.getIntensity() * 2.37f;

        byte regionByte;
        switch (payload.getTargetRegion()) {
            case OPTIC_NERVE -> regionByte = 0x01;
            case MOTOR_CORTEX -> regionByte = 0x02;
            case LIMBIC_SYSTEM -> regionByte = 0x04;
            default -> throw new IllegalArgumentException("> Invalid Target Region! " + payload.getTargetRegion());
        }

        byte[] commandBytes = payload.getCommandText().getBytes(StandardCharsets.UTF_8);
        byte[] rawData = new byte[1 + commandBytes.length];
        rawData[0] = regionByte;
        System.arraycopy(commandBytes, 0, rawData, 1, commandBytes.length);

        int statusCode = hardware.shockSystem(voltage, rawData);

        if (statusCode == -1) {
            throw new RejectionSyndromeException("> Short Circuit Error!");
        } else if (statusCode == 137) {
            throw new RejectionSyndromeException("> Neural Rejection Error!");
        } else if (statusCode != 0) {
            throw new RuntimeException("> Unknown Error!");
        }
    }
}
