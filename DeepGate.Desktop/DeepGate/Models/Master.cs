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

        public string GroupHeaderDateTime
        {
            get
            {
                DateTime yesterday = DateTime.Now.Date.AddDays(-1);

                if (DateTime.Date == yesterday)
                {
                    return "Yesterday";
                }

                if (DateTime.Date == DateTime.Now.Date)
                {
                    return "Today";
                }

                // For all other dates, return in the desired format
                return DateTime.ToString("dd MMM");
            }
        }

        public Master() { }
    }
}