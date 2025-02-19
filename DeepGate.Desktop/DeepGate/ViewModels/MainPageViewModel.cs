using System;
using System.Reflection.Metadata;
using CommunityToolkit.Mvvm.ComponentModel;
using DeepGate.Helpers;
using DeepGate.Interfaces;
using DeepGate.Models;
using DeepGate.Models;

namespace DeepGate.ViewModels;

public partial class MainPageViewModel: ObservableObject
{
    #region Fields

    private readonly IApiService apiService;
    private readonly IWallPaperService wallPaperService;

    #endregion

    #region Properties

    [ObservableProperty]
    public string? imageUrl;

    [ObservableProperty]
    public ChatCompletion chatCompletion = new ChatCompletion { Messages = new List<Message>(), Model = "" };


    #endregion

    public MainPageViewModel(IApiService apiService, IWallPaperService wallPaperService)
	{
        this.apiService = apiService;
        this.wallPaperService = wallPaperService;

        GetBackGroundImage();

        AddMockMessages();
	}

    private void AddMockMessages()
    {
        ChatCompletion.Messages = new List<Message>
        {
            new Message
            {
            Role = Constants.UserRole,
            Content = "Hello"
            },
            new Message
            {
            Role = Constants.SystemRole,
            Content = "Hi there! How can I assist you today?"
            },
            new Message
            {
            Role = Constants.UserRole,
            Content = "I'm curious about the Death Star project. Can you tell me more about it?"
            },
            new Message
            {
            Role = Constants.SystemRole,
            Content = "The Death Star project is a massive space station with enough power to destroy entire planets. It was designed to enforce the Emperor's rule through fear."
            },
            new Message
            {
            Role = Constants.UserRole,
            Content = "That sounds terrifying. How does it generate such immense power?"
            },
            new Message
            {
            Role = Constants.SystemRole,
            Content = "The Death Star's superlaser is powered by kyber crystals, which are capable of amplifying energy to unprecedented levels. This energy is then focused into a single, devastating beam."
            },
            new Message
            {
            Role = Constants.UserRole,
            Content = "Are there any weaknesses in its design?"
            },
            new Message
            {
            Role = Constants.SystemRole,
            Content = "There is a small thermal exhaust port, right below the main port. A precise hit from a proton torpedo could potentially trigger a chain reaction, leading to the destruction of the entire station."
            }
        };
    }

    private async void GetBackGroundImage()
    {
        // Ice 431xzv
        // Mountains 4gogwd
        ImageUrl = await wallPaperService.GetImageURLForId("4gogwd");
    }
}
