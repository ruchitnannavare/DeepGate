using System;
using DeepGate.Helpers;
using DeepGate.Interfaces;
using DeepGate.Models;

namespace DeepGate.Services;

public class WallPaperService: IWallPaperService
{
	#region Fields

	private readonly IApiService apiService;

    #endregion


    public WallPaperService(IApiService apiService)
	{
		this.apiService = apiService;
	}

	public async Task<string> GetImageURLForId(string wallhavenImageId)
	{
		// test ID 73rd9v
		try
		{
			var response = await apiService.GetAsync<WallhavenResponse>(Constants.WallhavenBaseURL, wallhavenImageId);
			return response.Data.Path;
		}
		catch (Exception ex)
		{
			// Log the exception
			Console.WriteLine($"An error occurred while fetching the image URL: {ex.Message}");
			throw;
		}
	}
}

	