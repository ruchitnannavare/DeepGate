using System.Collections.Generic;
using System.Threading.Tasks;
using DeepGate.Models;

namespace DeepGate.Interfaces;

public interface IDataBaseHelper
{
    Task<int> AddOrUpdateMasterInstance(Master masterModel);
    Task<List<Master>> GetAllInstances();

    Task<int> AddWallpaper(Wallpaper wallpaper);
    Task<List<Wallpaper>> GetAllWallpaper();

    Task SaveBackgroundSetting(BackgroundSetting setting);
    Task<BackgroundSetting> GetBackgroundSetting();
}