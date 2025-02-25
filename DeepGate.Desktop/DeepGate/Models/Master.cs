using System;
using CommunityToolkit.Mvvm.ComponentModel;
using SQLite;

namespace DeepGate.Models
{
	public partial class Master: ObservableObject
	{
        [PrimaryKey, AutoIncrement]
        public int Id { get; set; }

        [ObservableProperty]
        string summary;

        [ObservableProperty]
        DateTime dateTime;

        public Master() { }
    }
}

