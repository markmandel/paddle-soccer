using UnityEngine.Rendering;

public class PlayerInfo
{
    /// <summary>
    /// Detect headless mode (which has graphicsDeviceType Null)
    /// </summary>
    /// <returns>of the game is headless (has a gpu)</returns>
    public static bool IsHeadless()
    {
        return UnityEngine.SystemInfo.graphicsDeviceType == GraphicsDeviceType.Null;
    }
}