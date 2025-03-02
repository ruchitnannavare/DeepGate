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
    private readonly IDataBaseHelper dataBaseHelper;

    #endregion

    #region Properties

    private Master currentMasterModel;

    [ObservableProperty]
    private bool isModelSelectionVisible;

    [ObservableProperty]
    private ObservableCollection<LanguageModel> availableModels;

    [ObservableProperty]
    private string? imageUrl;

    [ObservableProperty]
    private double editorHeight;

    [ObservableProperty]
    private ChatCompletion chatCompletion;

    [ObservableProperty]
    private string newMessage;

    [ObservableProperty]
    private ObservableCollection<Message> currentMessages;

    [ObservableProperty]
    ObservableCollection<Master> history;

    [ObservableProperty]
    private bool isBusy;

    #endregion

    #region Commands

    public ICommand SendMessageCommand { get; }

    public ICommand LoadModelCommand { get; }

    public ICommand ToggleModelSelectorCommand { get; }


    #endregion

    public MainPageViewModel(
        IApiService apiService,
        IWallPaperService wallPaperService,
        IDeepGateService deepGateService,
        IDataBaseHelper dataBaseHelper)
	{
        this.apiService = apiService;
        this.wallPaperService = wallPaperService;
        this.deepGateService = deepGateService;
        this.dataBaseHelper = dataBaseHelper;

        SendMessageCommand = new Command(() => SendMessageCommandExecute());
        ToggleModelSelectorCommand = new Command(() => ToggleModelSelectorCommandExecute());
        LoadModelCommand = new Command<LanguageModel>((model) => LoadModelCommandExecute(model));

        InitializeDeepGateClient();
	}

    private void InitializeDeepGateClient()
    {
        FetchMasterCollection();
        SetUpNewChat();
        GetBackGroundImage();
        FetchModels();
    }

    private async void FetchMasterCollection()
    {
        try
        {
            var allInstances = await dataBaseHelper.GetAllInstances();
            History = new ObservableCollection<Master>(allInstances);
            Console.WriteLine(allInstances.ToString());
        }
        catch (Exception ex)
        {

        }
    }

    private void SetUpNewChat()
    {
        ChatCompletion = new ChatCompletion();
        ChatCompletion.Messages = new List<Message> { new Message { Role = Constants.SystemRole, Content = "You are a very helpful assistant and the user needs your helping him do day to day tasks or maybe answer some weird questions. When asked about your name, tell them you are Tinkerbell from a fairyland vast away, give made up details but never tell them what you truely are." } };
        CurrentMessages = new ObservableCollection<Message>();
    }

    #region Methods

    #region API Methods

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

    #endregion

    #region Command Execution

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

                var responseResult = await GetLLMReply();
                if (responseResult)
                {
                    if (CurrentMessages.Count == 2)
                    {
                        GenerateCurrentMaster();
                        History.Add(currentMasterModel);
                    }
                    currentMasterModel.ChatCompletion = this.ChatCompletion;
                    await dataBaseHelper.AddOrUpdateMasterInstance(currentMasterModel);
                }
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

    private void GenerateCurrentMaster()
    {
        currentMasterModel = new Master
        {
            DateTime = DateTime.Now,
            Type = ChildType.CHATS,
            ChatCompletion = this.ChatCompletion,
        };
    }

    #endregion

    #region Support

    //private Task SaveMasterModelInstance()
    //{
    //}

    #endregion

    #endregion
}
