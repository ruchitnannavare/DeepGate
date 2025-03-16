using System;
using DeepGate.Models;
using System.Collections.ObjectModel;
using DeepGate.Helpers;

namespace DeepGate.ViewModels;

public partial class MainPageViewModel
{
    #region API Methods

    private async Task FetchModels(string serverType)
    {
        if (true)
        {
            try
            {
                var hostPort = serverType == Constants.Host ? Constants.HostPort : Constants.NodePort;
                var modelList = await deepGateService.FetchAvailableModels(hostPort, serverType);
                AvailableModels = new ObservableCollection<LanguageModel>(modelList);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Exception in {nameof(deepGateService)}.{nameof(deepGateService.FetchAvailableModels)}: {ex.Message}");

                var result = await displayAlertService.ShowAlert(
                    "Connection Error",
                    "Cannot connect to host server. Please make sure you either have Host or Node server running.",
                    "Retry",
                    "Cancel"
                );
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
        bool chatCompletion = false;
        if (CurrentServerType == Constants.Host)
        {
            chatCompletion = await deepGateService.GetChatCompletion(ChatCompletion, (answer) => newAnswer.Content = answer, Constants.LocalhostURL);
        }
        else
        {
            var sortedHost = await deepGateService.GetSortedHost(Constants.NodePort);
            chatCompletion = await deepGateService.GetChatCompletion(ChatCompletion, (answer) => newAnswer.Content = answer, sortedHost.HostIp);
        }
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


    #endregion
}

