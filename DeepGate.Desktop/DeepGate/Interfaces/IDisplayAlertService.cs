using System;
namespace DeepGate.Interfaces
{
	public interface IDisplayAlertService
	{
		Task<string> ShowAlert(string title, string message, string optionA, string optionB);
	}
}

