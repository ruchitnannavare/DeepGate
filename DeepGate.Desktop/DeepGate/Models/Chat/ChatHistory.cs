namespace DeepGate.Models;

/// <summary>
/// Represents the history of a chat session.
/// </summary>
public class ChatHistory
{
    /// <summary>
    /// Gets or sets the unique identifier for the chat.
    /// </summary>
    public string ChatId { get; set; }

    /// <summary>
    /// Gets or sets the summary of the chat.
    /// </summary>
    public string Summary { get; set; }

    /// <summary>
    /// Gets or sets the chat completion details.
    /// </summary>
    public ChatCompletion ChatCompletion { get; set; }

    /// <summary>
    /// Gets the current model used in the chat completion.
    /// </summary>
    /// <returns>The current model as a string.</returns>
    public string CurrentModel()
    {
        return ChatCompletion.Model;
    }

    /// <summary>
    /// Changes the model used in the chat completion.
    /// </summary>
    /// <param name="chatHistory">The chat history object to update.</param>
    /// <param name="newModel">The new model to set.</param>
    public static void ChangeModel(ChatHistory chatHistory, string newModel)
    {
        chatHistory.ChatCompletion.Model = newModel;
    }
}