using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using DeepGate.Interfaces;

namespace DeepGate.Services;

public class ApiService : IApiService
{
    private readonly HttpClient httpClient;
    private readonly string apiKey;
    private readonly JsonSerializerOptions jsonOptions;

    public ApiService(HttpClient httpClient)
    {
        this.httpClient = httpClient;
        apiKey = "your_api_key_here"; // Replace with secure storage if needed

        // Set default JSON options to use camelCase
        jsonOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            DefaultIgnoreCondition = System.Text.Json.Serialization.JsonIgnoreCondition.WhenWritingNull
        };
    }

    public async Task<T> GetAsync<T>(string baseUrl, string endpoint)
    {
        var url = $"{baseUrl}/{endpoint}";
        //httpClient.DefaultRequestHeaders.Remove("Authorization");
        //httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {apiKey}");

        var response = await httpClient.GetAsync(url);
        return await HandleResponse<T>(response);
    }

    public async Task<T> PostAsync<T>(string baseUrl, string endpoint, object data)
    {
        var url = $"{baseUrl}/{endpoint}";
        //httpClient.DefaultRequestHeaders.Remove("Authorization");
        //httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {apiKey}");

        var jsonData = JsonSerializer.Serialize(data, jsonOptions);
        var content = new StringContent(jsonData, Encoding.UTF8, "application/json");

        var response = await httpClient.PostAsync(url, content);
        return await HandleResponse<T>(response);
    }

    private async Task<T> HandleResponse<T>(HttpResponseMessage response)
    {
        var json = await response.Content.ReadAsStringAsync();
        if (response.IsSuccessStatusCode)
        {
            return JsonSerializer.Deserialize<T>(json, jsonOptions)!;
        }
        else
        {
            throw new HttpRequestException($"Error: {response.StatusCode} - {json}");
        }
    }
}