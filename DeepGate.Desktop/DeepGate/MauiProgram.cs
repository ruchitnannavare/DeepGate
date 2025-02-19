using Microsoft.Extensions.Logging;
using Microsoft.Maui.LifecycleEvents;
using UIKit;
using ObjCRuntime;
using CoreGraphics;
using CommunityToolkit.Maui;
using CoreBTS.Maui.ShieldMVVM;
using CoreBTS.Maui.ShieldMVVM.Configuration;
using DeepGate.ViewModels;
using DeepGate.Interfaces;
using DeepGate.Services;
using Xe.AcrylicView;
using MauiIcons.Cupertino;
using MauiIcons.Material;
using MauiIcons.FontAwesome;
using MauiIcons.Fluent;
using Microsoft.Maui.Handlers;
using DeepGate.Views;
using Syncfusion.Maui.Core.Hosting;

namespace DeepGate;

public static class MauiProgram
{
    public static MauiApp CreateMauiApp()
	{
        var builder = MauiApp.CreateBuilder();
        builder
            .UseMauiApp<App>()
            .UseMauiCommunityToolkit()
            .UseAcrylicView()
            // Adding Icon dependencies
            .UseFluentMauiIcons()
            .UseFontAwesomeMauiIcons()
            .UseMaterialMauiIcons()
            .UseCupertinoMauiIcons()
            .ConfigureFonts(fonts =>
            {
                fonts.AddFont("OpenSans-Regular.ttf", "OpenSansRegular");
                fonts.AddFont("OpenSans-Semibold.ttf", "OpenSansSemibold");
            });

        ModifyEditor();
        builder.ConfigureSyncfusionCore();
        builder.ConfigureServices();

#if DEBUG
        builder.Logging.AddDebug();
#endif

        return builder.Build();
    }

    public static void ModifyEditor()
    {
        EditorHandler.Mapper.AppendToMapping("Borderless", (handler, view) =>
        {
#if MACCATALYST
            if (handler.PlatformView is UITextView textView)
            {
                textView.Layer.BorderWidth = 0;
                textView.BackgroundColor = UIColor.Clear;
            }
#endif
        });
    }

    private static MauiAppBuilder ConfigureServices(this MauiAppBuilder builder)
    {
        // ViewModels
        builder.Services.AddSingleton<MainPageViewModel>();

        // Views
        builder.Services.AddSingleton<MainPage>();

        // Services
        builder.Services.AddSingleton<HttpClient>();
        builder.Services.AddTransient<IApiService, ApiService>();
        builder.Services.AddTransient<IWallPaperService, WallPaperService>();

        return builder;
    }
}
