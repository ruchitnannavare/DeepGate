using DeepGate.ViewModels;
using MauiIcons.Core;

namespace DeepGate.Views;

public partial class MainPage : ContentPage
{
	public MainPage(MainPageViewModel mainPageViewModel)
	{
		InitializeComponent();
        // Temporary Workaround for url styled namespace in xaml
        _ = new MauiIcon();
        BindingContext = mainPageViewModel;
	}

}