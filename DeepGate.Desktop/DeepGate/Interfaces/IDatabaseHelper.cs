using System.Collections.Generic;
using System.Threading.Tasks;
using DeepGate.Models;

namespace DeepGate.Interfaces;

public interface IDataBaseHelper
{
    Task Init();
    Task<int> AddMasterInstance(Master masterModel);
    Task<List<Master>> AddMasterInstance();
}