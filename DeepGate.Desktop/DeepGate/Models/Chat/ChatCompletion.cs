using System;
using CommunityToolkit.Mvvm.ComponentModel;
using LiteDB;
using SQLite;

namespace DeepGate.Models;

/// <summary>
/// Represents a chat completion model.
/// </summary>
public partial class ChatCompletion: ObservableObject
{
	/// <summary>
	/// Gets or sets the model identifier.
	/// </summary>
	[BsonField]
	[ObservableProperty]
	public string? model;

	/// <summary>
	/// Gets or sets the list of messages.
	/// </summary>
	[BsonField]
	[ObservableProperty]
	public List<Message>? messages;

	public ChatCompletion() { }
}

