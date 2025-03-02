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

    public DatabaseHelper()
    {
        var dbPath = Path.Combine(FileSystem.AppDataDirectory, "deepgate_lite_test_9.db");

        BsonMapper.Global.Entity<Message>();
        BsonMapper.Global.Entity<ChatCompletion>();

        database = new LiteDatabase(dbPath);
        masterCollection = database.GetCollection<Master>("masters");
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
}