using System;
using CommunityToolkit.Mvvm.ComponentModel;
namespace DeepGate.Models;

public partial class Message: ObservableObject
{
	[ObservableProperty]
	public string role;

	[ObservableProperty]
	public string content;
}

