using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using DeepGate.Interfaces;
using Newtonsoft.Json;

namespace DeepGate.Services;

public class ApiService : IApiService
{
    public HttpClient HttpClient { get; }
    private readonly string apiKey;
    private readonly JsonSerializerOptions jsonOptions;

    public ApiService(HttpClient HttpClient)
    {
        this.HttpClient = HttpClient;
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
        //HttpClient.DefaultRequestHeaders.Remove("Authorization");
        //HttpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {apiKey}");

        var response = await HttpClient.GetAsync(url);
        return await HandleResponse<T>(response);
    }

    public async Task<T> PostAsync<T>(string baseUrl, string endpoint, object data)
    {
        var url = $"{baseUrl}/{endpoint}";
        //HttpClient.DefaultRequestHeaders.Remove("Authorization");
        //HttpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {apiKey}");

        var jsonData = JsonConvert.SerializeObject(data);
        var content = new StringContent(jsonData, Encoding.UTF8, "application/json");

        var response = await HttpClient.PostAsync(url, content);
        return await HandleResponse<T>(response);
    }

    private async Task<T> HandleResponse<T>(HttpResponseMessage response)
    {
        var json = await response.Content.ReadAsStringAsync();
        if (response.IsSuccessStatusCode)
        {
            return JsonConvert.DeserializeObject<T>(json)!;
        }
        else
        {
            throw new HttpRequestException($"Error: {response.StatusCode} - {json}");
        }
    }
}