using System.Threading.Tasks;

namespace DeepGate.Interfaces;

public interface IApiService
{
    HttpClient HttpClient { get; }
    Task<T> GetAsync<T>(string baseUrl, string endpoint);
    Task<T> PostAsync<T>(string baseUrl, string endpoint, object data);
}