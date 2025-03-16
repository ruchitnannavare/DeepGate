using System;
namespace DeepGate.Helpers;

public static class Constants
{
    public static string WallhavenBaseURL = "https://wallhaven.cc/api/v1/w";
    public const string UserRole = "user";
    public const string SystemRole = "system";
    public const string AssistantRole = "assistant";
    public const string Host = "host";
    public const string Node = "node";
    public static string LocalhostURL = "127.0.0.1";

    public static string BaseURLFormat = "http://{0}:{1}";

    public static string HostPort = "9090";
    public static string NodePort = "8080";
    public const string FirstBoot = "bootupSequence";


    // Background categories
    public static string Nature = "Nature";
    public static string Space = "Space";
    public static string Cityscape = "Cityscape";


    public static readonly List<(string, string)> SavedBackgrounds = new List<(string, string)>
    {
        // Nature
        ("4l8y6p", Nature),
        ("jxz2em", Nature),
        ("2kmpkx", Nature),
        ("4dyok3", Nature),
        ("nmowx1", Nature),

        // Space
        ("mdwg31", Space),
        ("j3zdom", Space),
        ("83ddr1", Space),
        ("0qgx6d", Space),
        ("d5j55o", Space),

        ("3ldy99", Cityscape),
        ("4ljoml", Cityscape),
        ("0p91m9", Cityscape),
        ("r7mkw7", Cityscape),
    };

    public static string BasePrompt = "You are a structured AI assistant. Provide complete and contextually correct responses. Format your output as a JSON object: {\"response\": \"Full, helpful answer.\", \"summary\": \"Brief topic summary.\"}. Ensure 'response' is well-formed and contains all necessary details before generating 'summary'. Never swap content between 'response' and 'summary'.";
}

