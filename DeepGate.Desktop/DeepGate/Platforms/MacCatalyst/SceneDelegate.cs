using UIKit;
using Foundation;
using CoreGraphics;

namespace DeepGate;

public class SceneDelegate : UIResponder, IUIWindowSceneDelegate
{
    public UIWindow? Window { get; set; } // Ensure Window is declared

    [Export("scene:willConnectToSession:options:")]
    public void SceneWillConnect(UIScene scene, UISceneSession session, UISceneConnectionOptions connectionOptions)
    {
        if (scene is UIWindowScene windowScene)
        {
            Window = new UIWindow(windowScene); // Set up UIWindow
            Window.RootViewController = new UIViewController(); // This should be replaced with your main app UI
            Window.MakeKeyAndVisible(); // Ensure it is displayed

            var restrictions = windowScene.SizeRestrictions;
            if (restrictions != null)
            {
                restrictions.MinimumSize = new CGSize(1000, 800);  // Set minimum window size
            }
        }
    }
}