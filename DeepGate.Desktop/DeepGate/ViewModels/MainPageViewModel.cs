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
    private readonly IDisplayAlertService displayAlertService;
    private readonly IDataBaseHelper dataBaseHelper;
    private readonly IPreferences preferences;

    private bool isFirstLaunch;

    #endregion

    #region Properties

    private Master currentMasterModel;

    [ObservableProperty]
    public string currentServerType;

    [ObservableProperty]
    private bool isModelSelectionVisible;

    [ObservableProperty]
    private ObservableCollection<LanguageModel> availableModels;

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

    public ICommand MasterSelectionCommand { get; }

    public ICommand NewChatCommand { get; }

    #endregion

    public MainPageViewModel(
        IApiService apiService,
        IWallPaperService wallPaperService,
        IDeepGateService deepGateService,
        IDataBaseHelper dataBaseHelper,
        IPreferences preferences,
        IDisplayAlertService displayAlertService)
    {
        this.apiService = apiService;
        this.wallPaperService = wallPaperService;
        this.deepGateService = deepGateService;
        this.dataBaseHelper = dataBaseHelper;
        this.preferences = preferences;
        this.displayAlertService = displayAlertService;

        isFirstLaunch = preferences.Get<string?>(Constants.FirstBoot, null) == null;

        SendMessageCommand = new Command(() => SendMessageCommandExecute());
        ToggleModelSelectorCommand = new Command(() => ToggleModelSelectorCommandExecute());
        LoadModelCommand = new Command<LanguageModel>((model) => LoadModelCommandExecute(model));
        MasterSelectionCommand = new Command<Master>((model) => MasterSelectionCommandExecute(model));
        NewChatCommand = new Command(() => NewChatCommandExecute());

        InitializeWallpaperCommands();

        InitializeDeepGateClient();
    }

    private void InitializeDeepGateClient()
    {
        FetchMasterCollection();
        SetUpNewChat();
        _ = ServerSelection();
    }

    private async Task ServerSelection()
    {
        var result = await displayAlertService.ShowAlert(
            "Server connection option",
            "Select what 0-898",
            Constants.Host,
            Constants.Node
        );
        switch (result)
        {
            case Constants.Host:
                CurrentServerType = Constants.Host;
                break;

            case Constants.Node:
                CurrentServerType = Constants.Node;
                break;

            default:
                await displayAlertService.ShowAlert(
                    "Error",
                    "Invalid server selection",
                    "OK",
                    "Cancel"
                );
                break;
        }

        _ = FetchModels(CurrentServerType);
    }

    private async void FetchMasterCollection()
    {
        var allInstances = await dataBaseHelper.GetAllInstances();
        History = new ObservableCollection<Master>(allInstances);
        Console.WriteLine(allInstances.ToString());
    }

    private void SetUpNewChat()
    {
        ChatCompletion = new ChatCompletion();
        ChatCompletion.Messages = new List<Message> { new Message { Role = Constants.SystemRole, Content = Constants.BasePrompt } };
        CurrentMessages = new ObservableCollection<Message>();
    }

    #region Methods

    

    #region Command Execution


    private void NewChatCommandExecute()
    {
        if (CurrentMessages.Count > 0)
        {
            NewMessage = string.Empty;
            SetUpNewChat();
        }
    }

    private void MasterSelectionCommandExecute(Master model)
    {
        NewMessage = string.Empty;
        ChatCompletion = model.ChatCompletion;
        currentMasterModel = model;
        if (ChatCompletion?.Messages != null && ChatCompletion.Messages.Count > 0)
        {
            CurrentMessages = new ObservableCollection<Message>(ChatCompletion.Messages.Skip(1));
        }
        else
        {
            CurrentMessages = new ObservableCollection<Message>();
        }
    }

    private async void LoadModelCommandExecute(LanguageModel model)
    {
        var hostPort = CurrentServerType == Constants.Host ? Constants.HostPort : Constants.NodePort;
        var modelLoaded = await deepGateService.LoadModel(model.Name, (status) => model.IsLoading = status, hostPort, CurrentServerType);

        if (modelLoaded)
        {
            ChatCompletion.Model = model.Name;
        }
    }

    private void ToggleModelSelectorCommandExecute()
    {
        IsModelSelectionVisible = !IsModelSelectionVisible;
    }

    private async void SendMessageCommandExecute()
    {
        if (!string.IsNullOrEmpty(NewMessage) & !string.IsNullOrEmpty(ChatCompletion.Model) & !IsBusy)
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
                        MainThread.BeginInvokeOnMainThread(() =>
                        {
                            History.ToList().Add(currentMasterModel);
                            History = new ObservableCollection<Master>(History.OrderByDescending(m => m.DateTime));
                        });
                    }
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

    #endregion

    #endregion
}
