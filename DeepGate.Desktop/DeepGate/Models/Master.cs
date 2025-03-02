using System;
using CommunityToolkit.Mvvm.ComponentModel;
using Newtonsoft.Json;
using LiteDB;

namespace DeepGate.Models
{
    public partial class Master : ObservableObject
    {
        [BsonId]
        public int Id { get; set; }

        [ObservableProperty]
        public ChatCompletion chatCompletion;

        /// <summary>
        /// Type of Master model
        /// <list type="0">CHATS</list>
        /// <list type="1">NOTES</list>
        /// <list type="2">PROMPTS</list>
        /// </summary>
        public ChildType Type { get; set; }

        [ObservableProperty]
        string summary;

        [ObservableProperty]
        DateTime dateTime;

        public Master() { }
    }
}