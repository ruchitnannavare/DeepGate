using SQLite;
using System.IO;
using System.Threading.Tasks;
using System.Collections.Generic;
using DeepGate.Interfaces;
using DeepGate.Models;

public class DatabaseHelper : IDataBaseHelper
{
    private readonly SQLiteAsyncConnection database;

    public DatabaseHelper(string dbPath)
    {
        database = new SQLiteAsyncConnection(dbPath);
    }

    public async Task Init()
    {
        await database.CreateTableAsync<Master>();
    }

    public async Task<int> AddMasterInstance(Master master)
    {
        return await database.InsertAsync(master);
    }

    public async Task<List<Master>> AddMasterInstance()
    {
        return await database.Table<Master>().ToListAsync();
    }
}