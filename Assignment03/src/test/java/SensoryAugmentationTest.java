package tests;

import internal.implementor.ImplantFirmware;
import internal.asbtraction.NeuralInterface;
import internal.implementor.NeuralPayload;
import internal.asbtraction.SensoryAugmentation;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

public class SensoryAugmentationTest {

    @Test
    void triggerImplant_delegatesCorrectPayloadToFirmware() {
        ImplantFirmware mockFirmware = mock(ImplantFirmware.class);
        NeuralInterface implant = new SensoryAugmentation(mockFirmware);
        implant.triggerImplant();
        ArgumentCaptor<NeuralPayload> captor = ArgumentCaptor.forClass(NeuralPayload.class);
        verify(mockFirmware, times(1)).executeNeuralCommand(captor.capture());

        NeuralPayload capturedPayload = captor.getValue();
        assertEquals(NeuralPayload.TargetRegion.OPTIC_NERVE, capturedPayload.getTargetRegion());
        assertEquals(3, capturedPayload.getIntensity());
        assertEquals("AUTH_ENC_NIGHT_VISION", capturedPayload.getCommandText());
    }
}
