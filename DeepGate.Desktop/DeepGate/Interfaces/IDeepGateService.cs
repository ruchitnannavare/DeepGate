using System;
using DeepGate.Models;

namespace DeepGate.Interfaces;

public interface IDeepGateService
{
    Task<bool> LoadModel(string modelName, Action<bool> loadingStatus, string serverEnvironment);

    Task<List<LanguageModel>> FetchAvailableModels(string serverEnvironment);

    Task<bool> GetChatCompletion(ChatCompletion chats, Action<string> answer, string serverEnvironment);
}
