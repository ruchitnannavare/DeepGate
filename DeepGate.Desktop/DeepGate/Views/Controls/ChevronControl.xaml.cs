using System.Windows.Input;

namespace DeepGate.Views;

public partial class ChevronControl : ContentView
{
	string unfocusedColor = "#30FFFFFF";
	string focusedColor = "#60FFFFFF";
	Action clickAction;

    public ChevronControl()
	{
		InitializeComponent();
	}

	private void OnLeftArrowClicked(System.Object sender, Microsoft.Maui.Controls.TappedEventArgs e)
	{
        // Handle left arrow click event
        LeftArrowCommand?.Execute(null);
	}

	private void OnRightArrowClicked(System.Object sender, Microsoft.Maui.Controls.TappedEventArgs e)
	{
		// Handle right arrow click event
		RightArrowCommand?.Execute(null);
	}

	public static readonly BindableProperty LeftArrowCommandProperty = BindableProperty.Create(
		nameof(LeftArrowCommand),
		typeof(ICommand),
		typeof(ChevronControl),
		default(ICommand));

	public ICommand LeftArrowCommand
	{
		get => (ICommand)GetValue(LeftArrowCommandProperty);
		set => SetValue(LeftArrowCommandProperty, value);
	}

	public static readonly BindableProperty RightArrowCommandProperty = BindableProperty.Create(
		nameof(RightArrowCommand),
		typeof(ICommand),
		typeof(ChevronControl),
		default(ICommand));

	public ICommand RightArrowCommand
	{
		get => (ICommand)GetValue(RightArrowCommandProperty);
		set => SetValue(RightArrowCommandProperty, value);
	}

	public static readonly BindableProperty SelectionTextProperty = BindableProperty.Create(
		nameof(SelectionText),
		typeof(string),
		typeof(ChevronControl),
		default(string),
		propertyChanged: OnSelectionTextChanged);

	public string SelectionText
	{
		get => (string)GetValue(SelectionTextProperty);
		set => SetValue(SelectionTextProperty, value);
	}

	private static void OnSelectionTextChanged(BindableObject bindable, object oldValue, object newValue)
	{
		var control = (ChevronControl)bindable;
		if (control.selectionLabel != null)
		{
			control.selectionLabel.Text = (string)newValue;
		}
	}
}
