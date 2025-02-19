using System;
namespace DeepGate.Interfaces;

public interface IWallPaperService
{
	Task<string> GetImageURLForId(string wallhavenImageId);
}


