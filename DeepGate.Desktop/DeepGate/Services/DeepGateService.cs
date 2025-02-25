using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using DeepGate.Helpers;
using DeepGate.Interfaces;
using DeepGate.Models;

namespace DeepGate.Services;


public class DeepGateService: IDeepGateService
{
    #region Fields

    private readonly IApiService apiService;

    #endregion

    public const string LoadModelApi = "/load-model";
    public const string ChatCompletionApi = "/chat";
    public const string FetchModelListApi = "/fetch-models";

    public DeepGateService(IApiService apiService)
    {
        this.apiService = apiService;
    }

    public async Task<bool> LoadModel(string modelName, Action<bool> loadingStatus, string serverEnvironment = Constants.Host)
    {
        if (string.IsNullOrEmpty(modelName))
        {
            throw new ArgumentException("Model name cannot be empty", nameof(modelName));
        }

        try
        {
            //ModelLoadingStatus = $"Loading model: {modelName}";
            loadingStatus(true);
            var baseURL = Constants.LocalhostBaseURL + Constants.HostPort;
            var endPoint = serverEnvironment + LoadModelApi;

            var payload = new Dictionary<string, object>
            {
                { "model_name", modelName }
            };

            var response = await apiService.PostAsync<Dictionary<string, object>>(baseURL, endPoint, payload);

            // If we get here, the model loaded successfully
            //ModelLoadingStatus = $"Successfully loaded model: {modelName}";
            return true;
        }
        catch (Exception ex)
        {
            //ModelLoadingStatus = $"Error loading model {modelName}: {ex.Message}";
            throw new Exception($"Failed to load model {modelName}", ex);
        }
        finally
        {
            loadingStatus(false);
        }
    }

    public async Task<List<LanguageModel>> FetchAvailableModels(string serverEnvironment = Constants.Host)
    {
        try
        {
            var baseURL = Constants.LocalhostBaseURL + Constants.HostPort;
            var endPoint = serverEnvironment + FetchModelListApi;

            var models = await apiService.GetAsync<ModelsList>(baseURL, endPoint);
            return models.Models;
        }
        catch (Exception ex)
        {
            // Log the error or handle it as needed
            throw new Exception($"Failed to fetch models: {ex.Message}", ex);
        }
    }

    public async Task<bool> GetChatCompletion(ChatCompletion chats, Action<string> answer, string serverEnvironment = Constants.Host)
    {
        try
        {
            var baseURL = Constants.LocalhostBaseURL + Constants.HostPort;
            var endPoint = serverEnvironment + ChatCompletionApi;

            using var request = new HttpRequestMessage(HttpMethod.Post, $"{baseURL}/{endPoint}");
            request.Headers.Accept.Add(new System.Net.Http.Headers.MediaTypeWithQualityHeaderValue("text/event-stream"));

            var jsonData = JsonSerializer.Serialize(chats);
            request.Content = new StringContent(jsonData, Encoding.UTF8, "application/json");

            // Use the apiService's httpClient directly for streaming
            var response = await apiService.HttpClient.SendAsync(
                request,
                HttpCompletionOption.ResponseHeadersRead);

            response.EnsureSuccessStatusCode();

            using var stream = await response.Content.ReadAsStreamAsync();
            using var reader = new StreamReader(stream);

            string messageBuilder = "";
            bool thinkTag = false;

            // Run streaming in a background task
            await Task.Run(async () =>
            {
                while (!reader.EndOfStream)
                {
                    var line = await reader.ReadLineAsync();

                    if (string.IsNullOrWhiteSpace(line)) continue; // Skip empty lines

                    Console.WriteLine($"Received line: {line}");

                    if (line.StartsWith("data:", StringComparison.OrdinalIgnoreCase))
                    {
                        var data = line.Substring(5);
                        if (data.Equals("[DONE]", StringComparison.OrdinalIgnoreCase))
                            break;

                        if (data.Equals("<think>"))
                        {
                            thinkTag = true;
                        }

                        if (!thinkTag)
                        {
                            messageBuilder += data;
                            MainThread.BeginInvokeOnMainThread(() =>
                            {
                                answer(messageBuilder); // Update UI without freezing
                            });
                        }

                        if (data.Equals("</think>"))
                        {
                            thinkTag = false;
                        }
                    }
                }
            });
            return true;
        }
        catch (Exception ex)
        {
            System.Console.WriteLine("An error occurred while fetching chat completion: " + ex.Message);
            return false;
        }
    }
}