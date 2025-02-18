using Microsoft.Extensions.Logging;
using Microsoft.Maui.LifecycleEvents;
using UIKit;
using ObjCRuntime;
using CoreGraphics;
using CommunityToolkit.Maui;
using CoreBTS.Maui.ShieldMVVM;
using CoreBTS.Maui.ShieldMVVM.Configuration;
using DeepGate.ViewModels;

namespace DeepGate;

public static class MauiProgram
{
    public static MauiApp CreateMauiApp()
	{
        var builder = MauiApp.CreateBuilder();
        builder
            .UseMauiApp<App>()
            .UseMauiCommunityToolkit()
            .ConfigureFonts(fonts =>
            {
                fonts.AddFont("OpenSans-Regular.ttf", "OpenSansRegular");
                fonts.AddFont("OpenSans-Semibold.ttf", "OpenSansSemibold");
            });

#if DEBUG
        builder.Logging.AddDebug();
#endif

        return builder.Build();
    }

    private static MauiAppBuilder ConfigureServices(this MauiAppBuilder builder)
    {
        builder.Services.AddSingleton<MainPageViewModel>();
        builder.Services.AddSingleton<MainPage>();

        return builder;
    }
}
