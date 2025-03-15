using LiteDB;
using System.IO;
using System.Collections.Generic;
using DeepGate.Interfaces;
using DeepGate.Models;
using Microsoft.Maui.Storage; // For cross-platform file path


public class DatabaseHelper : IDataBaseHelper
{
    private readonly LiteDatabase database;
    private readonly ILiteCollection<Master> masterCollection;
    private readonly ILiteCollection<Wallpaper> wallpaperCollection;
    private readonly ILiteCollection<BackgroundSetting> backgroundSettingCollection;

    public DatabaseHelper()
    {
        var dbPath = Path.Combine(FileSystem.AppDataDirectory, "deepgate_lite_test_9.db");

        BsonMapper.Global.Entity<Message>();
        BsonMapper.Global.Entity<ChatCompletion>();
        BsonMapper.Global.Entity<WallhavenData>();
        BsonMapper.Global.Entity<Wallpaper>();

        database = new LiteDatabase(dbPath);
        masterCollection = database.GetCollection<Master>("masters");
        wallpaperCollection = database.GetCollection<Wallpaper>("wallpapers");
        backgroundSettingCollection = database.GetCollection<BackgroundSetting>("backgroundsetting");
    }

    public Task<int> AddOrUpdateMasterInstance(Master master)
    {
        masterCollection.Upsert(master);
        return Task.FromResult(1); // Upsert doesn't return row count, assuming success
    }

    public Task<List<Master>> GetAllInstances()
    {
        var result = masterCollection.FindAll().ToList();
        return Task.FromResult(result);
    }

    public Task<int> AddWallpaper(Wallpaper wallpaper)
    {
        wallpaperCollection.Upsert(wallpaper);
        return Task.FromResult(1); // Upsert doesn't return row count, assuming success
    }

    public Task<List<Wallpaper>> GetAllWallpaper()
    {
        var result = wallpaperCollection.FindAll().ToList();
        return Task.FromResult(result);
    }

    public Task SaveBackgroundSetting(BackgroundSetting setting)
    {
        // Always use ID 1 for the single instance
        setting.Id = 1;
        backgroundSettingCollection.Upsert(setting);
        return Task.CompletedTask;
    }

    public Task<BackgroundSetting> GetBackgroundSetting()
    {
        // Retrieve the single instance with ID 1
        var setting = backgroundSettingCollection.FindById(1);
        return Task.FromResult(setting);
    }
}