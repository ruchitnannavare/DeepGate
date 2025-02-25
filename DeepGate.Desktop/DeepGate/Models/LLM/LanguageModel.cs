using System;
using CommunityToolkit.Mvvm.ComponentModel;
using Newtonsoft.Json;

namespace DeepGate.Models;

public class ModelsList
{
    public List<LanguageModel> Models { get; set; }
}

public partial class LanguageModel: ObservableObject
{
    public ModelDetails Details { get; set; }
    public string Model { get; set; }
    public string Name { get; set; }

    //View Properties
    [ObservableProperty]
    public bool isRunning;

    [ObservableProperty]
    public bool isLoading;
}

public class ModelDetails
{
    public string Family { get; set; }

    [JsonProperty("parameter_size")]
    public string ParameterSize { get; set; }  // Fix property name
}