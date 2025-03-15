using System;
using DeepGate.Models;

namespace DeepGate.Interfaces;

public interface IWallPaperService
{
	Task<WallhavenResponse> GetImageURLForId(string wallhavenImageId);
}


