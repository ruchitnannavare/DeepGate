using System.Collections.Generic;
using System.Threading.Tasks;
using DeepGate.Models;

namespace DeepGate.Interfaces;

public interface IDataBaseHelper
{
    Task<int> AddOrUpdateMasterInstance(Master masterModel);
    Task<List<Master>> GetAllInstances();
}