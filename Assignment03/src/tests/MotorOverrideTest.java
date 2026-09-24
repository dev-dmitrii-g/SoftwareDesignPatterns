package tests;

import internal.asbtraction.MotorOverride;
import internal.asbtraction.NeuralInterface;
import internal.implementor.NeuralPayload;
import internal.implementor.ImplantFirmware;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

public class MotorOverrideTest {

    @Test
    void triggerImplant_delegatesHighIntensityPayloadToFirmware() {
        ImplantFirmware mockFirmware = mock(ImplantFirmware.class);
        NeuralInterface implant = new MotorOverride(mockFirmware);
        implant.triggerImplant();
        ArgumentCaptor<NeuralPayload> captor = ArgumentCaptor.forClass(NeuralPayload.class);
        verify(mockFirmware, times(1)).executeNeuralCommand(captor.capture());

        NeuralPayload capturedPayload = captor.getValue();
        assertEquals(NeuralPayload.TargetRegion.MOTOR_CORTEX, capturedPayload.getTargetRegion());
        assertEquals(10, capturedPayload.getIntensity());
    }
}
