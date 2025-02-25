using System;
using System.Collections.ObjectModel;
using System.Reflection.Metadata;
using System.Windows.Input;
using CommunityToolkit.Maui.Core.Platform;
using CommunityToolkit.Mvvm.ComponentModel;
using DeepGate.Helpers;
using DeepGate.Interfaces;
using DeepGate.Models;

namespace DeepGate.ViewModels;

public partial class MainPageViewModel: ObservableObject
{
    #region Fields

    private readonly IApiService apiService;
    private readonly IWallPaperService wallPaperService;
    private readonly IDeepGateService deepGateService;

    #endregion

    #region Properties

    [ObservableProperty]
    private bool isModelSelectionVisible;

    [ObservableProperty]
    private ObservableCollection<LanguageModel> availableModels;

    [ObservableProperty]
    private string? imageUrl;

    [ObservableProperty]
    private double editorHeight;

    [ObservableProperty]
    private ChatCompletion chatCompletion = new ChatCompletion();

    [ObservableProperty]
    private string newMessage;

    [ObservableProperty]
    private ObservableCollection<Message> currentMessages;

    [ObservableProperty]
    private bool isBusy;

    #endregion

    #region Commands

    public ICommand SendMessageCommand { get; }

    public ICommand LoadModelCommand { get; }

    public ICommand ToggleModelSelectorCommand { get; }


    #endregion

    public MainPageViewModel(IApiService apiService, IWallPaperService wallPaperService, IDeepGateService deepGateService)
	{
        this.apiService = apiService;
        this.wallPaperService = wallPaperService;
        this.deepGateService = deepGateService;

        SendMessageCommand = new Command(() => SendMessageCommandExecute());
        ToggleModelSelectorCommand = new Command(() => ToggleModelSelectorCommandExecute());
        LoadModelCommand = new Command<LanguageModel>((model) => LoadModelCommandExecute(model));

        InitializeDeepGateClient();
	}

    private void InitializeDeepGateClient()
    {
        SetUpNewChat();
        GetBackGroundImage();
        FetchModels();
    }

    private void SetUpNewChat()
    {

#if MACCATALYST
        var content = new UserNotifications.UNMutableNotificationContent
        {
            Title = "Warning! ",
            Body = "This is an alert!",
            CategoryIdentifier = "WARNING_ALERT",
            Sound = UserNotifications.UNNotificationSound.DefaultCriticalSound
        };


        var trigger = UserNotifications.UNTimeIntervalNotificationTrigger.CreateTrigger(5, false);


        var request = UserNotifications.UNNotificationRequest.FromIdentifier("ALERT_REQUEST", content, trigger);


        UserNotifications.UNUserNotificationCenter.Current.AddNotificationRequest(request, (error) =>
        {
            if (error != null)
            {
                Console.WriteLine("Error: " + error);
            }
        });
#endif
        ChatCompletion.Messages = new List<Message> { new Message { Role = Constants.SystemRole, Content = "You are a very helpful assistant and the user needs your helping him do day to day tasks or maybe answer some weird questions. When asked about your name, tell them you are Tinkerbell from a fairyland vast away, give made up details but never tell them what you truely are." } };
        CurrentMessages = new ObservableCollection<Message>();
    }

    private async void LoadModelCommandExecute(LanguageModel model)
    {
        var modelLoaded = await deepGateService.LoadModel(model.Name, (status) => model.IsLoading = status, Constants.Host);
        if (modelLoaded)
        {
            ChatCompletion.Model = model.Name;
        }
    }

    private void ToggleModelSelectorCommandExecute()
    {
        IsModelSelectionVisible = !IsModelSelectionVisible;
    }

    private async void FetchModels()
    {
        try
        {
            var modelList = await deepGateService.FetchAvailableModels(Constants.Host);
            AvailableModels = new ObservableCollection<LanguageModel>(modelList);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Exception in {nameof(deepGateService)}.{nameof(deepGateService.FetchAvailableModels)}: {ex.Message}");

            // Show an alert in the UI
            await MainThread.InvokeOnMainThreadAsync(async () =>
            {
                bool retry = await Application.Current.MainPage.DisplayAlert(
                    "Connection Error",
                    "Cannot connect to host server. Please make sure you either have Host or Node server running.",
                    "Retry",
                    "Cancel"
                );

                if (retry)
                {
                    FetchModels(); // Retry fetching
                }
            });
        }
    }

    private async void SendMessageCommandExecute()    {
        if (!string.IsNullOrEmpty(NewMessage) & !string.IsNullOrEmpty(ChatCompletion.Model))
        {
            try
            {
                IsBusy = true;

                var newQuery = new Message
                {
                    Role = Constants.UserRole,
                    Content = NewMessage,
                };

                ChatCompletion?.Messages?.Add(newQuery);
                CurrentMessages.Add(newQuery);

                NewMessage = string.Empty;

                await GetLLMReply();

            }
            catch (Exception ex)
            {
                Console.WriteLine($"Exception in {nameof(MainPageViewModel)}.{nameof(SendMessageCommandExecute)}: {ex.Message}");
            }
            finally
            {
                IsBusy = false;
            }
        }
    }

    private async Task<bool> GetLLMReply()
    {
        var newAnswer = new Message
        {
            Role = Constants.AssistantRole,
            Content = "",
        };

        CurrentMessages.Add(newAnswer);
        var chatCompletion = await deepGateService.GetChatCompletion(ChatCompletion, (answer) => newAnswer.Content = answer, Constants.Host);
        if (!chatCompletion)
        {
            await MainThread.InvokeOnMainThreadAsync(async () =>
            {
                bool retry = await Application.Current.MainPage.DisplayAlert(
                    "Connection Error",
                    "Cannot connect to host server. Please make sure you either have Host or Node server running.",
                    "Retry",
                    "Cancel"
                );

                if (retry)
                {
                    await GetLLMReply(); // Retry getting LLM reply
                }
            });
        }

        //TODO: Add return value optimization later
        newAnswer.IsCompleted = true;
        ChatCompletion?.Messages?.Add(newAnswer);
        return true;
    }

    private async void GetBackGroundImage()
    {
        // Ice 431xzv
        // Mountains 4gogwd
        ImageUrl = await wallPaperService.GetImageURLForId("4gogwd");
    }
}
