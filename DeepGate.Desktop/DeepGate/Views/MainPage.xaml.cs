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

	public static readonly BindableProperty ModelSelectorVisibilityProperty =
		BindableProperty.Create(
			nameof(ModelSelectorVisibility),
			typeof(bool),
			typeof(MainPage),
			false,
			propertyChanged: OnModelSelectorVisibilityChanged);

	public bool ModelSelectorVisibility
	{
		get => (bool)GetValue(ModelSelectorVisibilityProperty);
		set => SetValue(ModelSelectorVisibilityProperty, value);
	}

	private static void OnModelSelectorVisibilityChanged(BindableObject bindable, object oldValue, object newValue)
	{
		var control = (MainPage)bindable;
		bool isVisible = (bool)newValue;
		if(isVisible)
		{
			//control.modelSelector.IsVisible = true;
			control.modelSelector.FadeTo(1, 350);
		}
		else
		{
			control.modelSelector.FadeTo(0, 350);
			//control.modelSelector.IsVisible = false;
		}
	}
}