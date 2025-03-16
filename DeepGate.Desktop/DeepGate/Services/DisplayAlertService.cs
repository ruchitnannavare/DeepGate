using System;
using Microsoft.Maui.Controls;
using DeepGate.Interfaces;

namespace DeepGate.Services;

public class DisplayAlertService : IDisplayAlertService
{
	public DisplayAlertService()
	{
	}

	public async Task<string> ShowAlert(string title, string message, string optionA, string optionB)
	{
		bool result = await MainThread.InvokeOnMainThreadAsync(async () =>
		{
			return await Application.Current.MainPage.DisplayAlert(
				title,
				message,
				optionA,
				optionB
			);
		});

		return result ? optionA : optionB;
	}
}

