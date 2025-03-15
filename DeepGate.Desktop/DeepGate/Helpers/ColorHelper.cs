using System;
namespace DeepGate.Helpers;

public static class ColorHelper
{
    public static string GetOpaqueColor(int opacity, string colorHex)
    {
        var opaqueColor = $"#{opacity}{colorHex}";
        return opaqueColor;
    }
}

