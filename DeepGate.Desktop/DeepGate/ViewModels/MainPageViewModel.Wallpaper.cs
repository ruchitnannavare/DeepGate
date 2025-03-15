using System;
using System.Windows.Input;
using CommunityToolkit.Mvvm.ComponentModel;
using DeepGate.Helpers;
using DeepGate.Models;

namespace DeepGate.ViewModels;

public partial class MainPageViewModel
{
    #region View Properties

    [ObservableProperty]
    string backgroundTintColor;

    [ObservableProperty]
    string backgroundAvatar;

    public List<Wallpaper> onDeviceBackgorunds = new List<Wallpaper>();

    public List<Wallpaper> wallpaperCollectionList = new List<Wallpaper>();

    public List<Wallpaper> selectedWallpaperCollectionList = new List<Wallpaper>();

    [ObservableProperty]
    private string? imageUrl;

    [ObservableProperty]
    private string? opacitySelectionText;

    [ObservableProperty]
    private string? wallpaperSelectionText = Constants.Cityscape;

    [ObservableProperty]
    private Wallpaper? selectedCollection;

    [ObservableProperty]
    private Wallpaper? selectedWallpaper;

    [ObservableProperty]
    int currentOpacity;

    BackgroundSetting currentBackgroundSetting;

    int wallpaperCollectionIndex= 1;

    #endregion



    #region Commands

    public ICommand LeftClickWallpaper { get; private set; }

    public ICommand RightClickWallpaper { get; private set; }

    public ICommand LeftClickOpacity { get; private set; }

    public ICommand RightClickOpacity { get; private set; }

    public ICommand ShuffleCommand { get; private set; }

    public ICommand SelectCollectionCommand { get; private set; }


    #endregion

    #region Methods

    private void InitializeWallpaperCommands()
    {
		// Initialize LeftClickWallpaper command
		LeftClickWallpaper = new Command(() => LeftClickWallpaperExecute());

		// Initialize RightClickWallpaper command
		RightClickWallpaper = new Command(() => RightClickWallpaperExecute());

		// Initialize LeftClickOpacity command
		LeftClickOpacity = new Command(() => LeftClickOpacityExecute());

		// Initialize RightClickOpacity command
		RightClickOpacity = new Command(() => RightClickOpacityExecute());

        // Initialize ShuffleCommand
        ShuffleCommand = new Command(() => ShuffleCommandExecute());

        // Initialize SelectCollectionCommand
        SelectCollectionCommand = new Command(() => SelectCollectionCommandExecute());

        FetchLocalBackgrounds(isFirstLaunch);
    }

    private async void SetBackgroundSettings(bool isFirstLaunch)
    {
        if (isFirstLaunch)
        {
            SetSelectedCollection(wallpaperCollectionIndex, wallpaperCollectionList);
            SelectedWallpaper = selectedWallpaperCollectionList.FirstOrDefault();
            SetOpacity(currentOpacity = 60);
            currentBackgroundSetting =
            new BackgroundSetting
            {
                SelectedWallpaper = this.SelectedWallpaper,
                SelectedWallpaperIndex = wallpaperCollectionIndex,
                Opacity = currentOpacity
            };
            await dataBaseHelper.SaveBackgroundSetting(currentBackgroundSetting);
        }
        else
        {
            SetSelectedCollection(wallpaperCollectionIndex, wallpaperCollectionList);
            SelectedWallpaper = selectedWallpaperCollectionList.FirstOrDefault();
            SetOpacity(60);
            //var savedBackground = currentBackgroundSetting = await dataBaseHelper.GetBackgroundSetting();
            //SetSelectedCollection(savedBackground.SelectedWallpaperIndex, wallpaperCollectionList);
            //SelectedWallpaper = savedBackground.SelectedWallpaper;
            //SetOpacity(savedBackground.Opacity);
        }
    }

    private async void FetchLocalBackgrounds(bool isFirstLaunch)
    {
        // Action to save wallpaers in db and current list
        Action<Wallpaper> addWallAction = new Action<Wallpaper>((wall) =>
        {
            onDeviceBackgorunds.Add(wall);
            dataBaseHelper.AddWallpaper(wall);
        });
        if (isFirstLaunch)
        {
            // Fire tasks to fetch saved backgrounds all at once
            foreach (var wallpaperId in Constants.SavedBackgrounds)
            {
                await AddWallpaperTask(wallpaperId.Item1, wallpaperId.Item2, addWallAction);
            }

            preferences.Set<string>(Constants.FirstBoot, "Complete");
        }
        else
        {
            onDeviceBackgorunds = await dataBaseHelper.GetAllWallpaper();
        }

        wallpaperCollectionList = onDeviceBackgorunds
        .GroupBy(w => w.CollectionName)
        .Select(g => g.First())
        .ToList();

        SetBackgroundSettings(isFirstLaunch);
    }

    private async Task AddWallpaperTask(string wallhavenId, string collectionName, Action<Wallpaper> listAdditionAction)
    {
        var wallpaper = new Wallpaper
        {
            WallHavenId = wallhavenId,
            CollectionName = collectionName,
        };

        wallpaper.Response = await wallPaperService.GetImageURLForId(wallhavenId);

        listAdditionAction(wallpaper);
    }

    private async void RightClickOpacityExecute()
    {
        if (currentOpacity < 90)
        {
            currentOpacity += 10;
            SetOpacity(currentOpacity);
            currentBackgroundSetting.Opacity = currentOpacity;
            await SaveBackgroundSetting(currentBackgroundSetting);
        }
    }

    private async void LeftClickOpacityExecute()
    {
        if (currentOpacity > 0)
        {
            currentOpacity -= 10;
            SetOpacity(currentOpacity);
            currentBackgroundSetting.Opacity = currentOpacity;
            await SaveBackgroundSetting(currentBackgroundSetting);
        }
    }

    private void RightClickWallpaperExecute()
    {
        if (wallpaperCollectionIndex < wallpaperCollectionList.Count)
        {
            wallpaperCollectionIndex += 1;
            SetSelectedCollection(wallpaperCollectionIndex, wallpaperCollectionList);
        }
    }

    private void LeftClickWallpaperExecute()
    {
        if (wallpaperCollectionIndex > 0)
        {
            wallpaperCollectionIndex -= 1;
            SetSelectedCollection(wallpaperCollectionIndex, wallpaperCollectionList);
        }
    }

    private void ShuffleCommandExecute()
    {
        if (selectedWallpaperCollectionList != null && selectedWallpaperCollectionList.Count > 0)
        {
            var random = new Random();
            var randomWallpaper = selectedWallpaperCollectionList[random.Next(selectedWallpaperCollectionList.Count)];
            
            if (!randomWallpaper.Equals(SelectedWallpaper))
            {
                SelectedWallpaper = randomWallpaper;
            }
        }
    }

    private async void SelectCollectionCommandExecute()
    {
        if (wallpaperCollectionIndex != currentBackgroundSetting.SelectedWallpaperIndex)
        {
            currentBackgroundSetting.SelectedWallpaperIndex = wallpaperCollectionIndex;
            await SaveBackgroundSetting(currentBackgroundSetting);
        }
    }

    private async Task SaveBackgroundSetting(BackgroundSetting setting)
    {
        await dataBaseHelper.SaveBackgroundSetting(setting);
    }

    private void SetOpacity(int currentOpacity)
    {
        BackgroundTintColor = ColorHelper.GetOpaqueColor(currentOpacity, "000000");
        OpacitySelectionText = $"{currentOpacity} %";
    }

    private void SetSelectedCollection(int index, List<Wallpaper> wallpapers)
    {
        SelectedCollection = wallpapers[index];
        selectedWallpaperCollectionList = onDeviceBackgorunds
        .Where(w => w.CollectionName == SelectedCollection.CollectionName)
        .ToList();
    }

    #endregion
}