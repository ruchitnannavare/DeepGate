using ConciergeGenie.Helpers;
using DeepGate.Views;
using UIKit;

namespace DeepGate;

public partial class App : Application
{
	public App(IServiceProvider serviceProvider)
	{
		InitializeComponent();
        // This licence key is required for including syncfusion controls The
        // key is added by the build pipeline
        Syncfusion.Licensing.SyncfusionLicenseProvider.RegisterLicense(ApiKeys.SyncfusionLicenceKey);
        MainPage = serviceProvider.GetService<MainPage>();
	}

    //protected override Window CreateWindow(IActivationState activationState)
    //{
    //    var window = new Window();
    //    window.HandlerChanged += windowHandlerChanged;
    //    var rootPage = new MainPage(); // Or whatever your root page is
    //    window.Page = rootPage;
    //    return window;
    //}

    //private void windowHandlerChanged(object sender, EventArgs e)
    //{
    //    var win = sender as Microsoft.Maui.Controls.Window;
    //    var uiWin = win.Handler.PlatformView as UIWindow;
    //    if (uiWin != null)
    //    {
    //        uiWin.WindowScene.Titlebar.TitleVisibility = UITitlebarTitleVisibility.Hidden;
    //    }
    //}
}
