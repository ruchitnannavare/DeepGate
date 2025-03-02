using DeepGate.Models;
using LiteDB;
using System.Collections.Generic;

namespace DeepGate.Helpers;

public class CustomBsonMapper : BsonMapper
{
    public CustomBsonMapper()
    {
        // Register Message class
        Global.RegisterType<Message>(
            serialize: msg => new BsonDocument
            {
                ["role"] = msg.Role,
                ["content"] = msg.Content,
                ["isCompleted"] = msg.IsCompleted
            },
            deserialize: bson =>
                new Message
                {
                    Role = bson["role"].AsString,
                    Content = bson["content"].AsString,
                    IsCompleted = bson["isCompleted"].AsBoolean
                }
        );

        // Register ChatCompletion class
        Global.RegisterType<ChatCompletion>(
            serialize: chat => new BsonDocument
            {
                ["model"] = chat.Model,
                ["messages"] = new BsonArray(chat.Messages?.ConvertAll(m => this.Serialize(m)) ?? new List<BsonValue>())
            },
            deserialize: bson =>
            new ChatCompletion
            {
                Model = bson["model"].AsString,
                Messages = bson["messages"].AsArray?.Select(m => this.Deserialize<Message>(m)).ToList() ?? new List<Message>()
            }
        );
    }
}