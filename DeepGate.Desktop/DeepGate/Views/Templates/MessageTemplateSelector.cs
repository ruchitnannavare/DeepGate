
using System.Collections.Immutable;
using DeepGate.Helpers;
using DeepGate.Models;
using Microsoft.Maui.Controls;

namespace DeepGate.Views;

[ContentProperty("Content")]
internal class MessageDataTemplateSelector : DataTemplateSelector
{
    public DataTemplate UserMessageTemplate { get; set; }
    public DataTemplate ReceivedMessageTemplate { get; set; }
    public DataTemplate ThinkingMessageTemplate { get; set; }

    protected override DataTemplate OnSelectTemplate(object item, BindableObject container)
    {
        var message = (Message)item;

        switch (message.Role)
        {
            case "user": 
                return UserMessageTemplate;
            case "assistant":
                return ReceivedMessageTemplate;
            default:
                return new DataTemplate();
        }
    }
}