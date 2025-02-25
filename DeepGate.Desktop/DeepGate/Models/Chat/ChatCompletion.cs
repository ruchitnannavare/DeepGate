using System;
using CommunityToolkit.Mvvm.ComponentModel;
namespace DeepGate.Models;

/// <summary>
/// Represents a chat completion model.
/// </summary>
public partial class ChatCompletion: Master
{
	/// <summary>
	/// Gets or sets the model identifier.
	/// </summary>
	[ObservableProperty]
	public string? model;

	/// <summary>
	/// Gets or sets the list of messages.
	/// </summary>
	[ObservableProperty]
	public List<Message>? messages;
}

