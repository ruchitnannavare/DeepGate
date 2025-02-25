using System;
using CommunityToolkit.Mvvm.ComponentModel;
namespace DeepGate.Models;

/// <summary>
/// Represents a chat message with a role, content, and completion status.
/// </summary>
public partial class Message : ObservableObject
{
	/// <summary>
	/// Gets or sets the role of the message sender.
	/// </summary>
	[ObservableProperty]
	public string role;

	/// <summary>
	/// Gets or sets the content of the message.
	/// </summary>
	[ObservableProperty]
	public string content;

	/// <summary>
	/// Gets or sets a value indicating whether the message is completed. For System ROle only.
	/// </summary>
	public bool IsCompleted { get; set; }
}

