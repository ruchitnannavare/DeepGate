using System;
using DeepGate.Models;
using DeepGate.Models.Responses;

namespace DeepGate.Interfaces;

public interface IDeepGateService
{
    Task<bool> LoadModel(string modelName, Action<bool> loadingStatus, string port, string serverEnvironment);

    Task<List<LanguageModel>> FetchAvailableModels(string port, string serverEnvironment);

    Task<bool> GetChatCompletion(ChatCompletion chats, Action<string> answer, string hostIp);

    Task<NodeSortedHostResponse> GetSortedHost(string port);
}
