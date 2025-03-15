using DeepGate.Models;
using Syncfusion.Maui.DataSource.Extensions;
namespace DeepGate.Helpers;

public class MasterGroupComparer : IComparer<GroupResult>
{
    public int Compare(GroupResult x, GroupResult y)
    {
        var xGroupItem = x.Items.Cast<Master>; 
        var yGroupItem = y.Items.Cast<Master>;
        var xObject = new List<Master>(xGroupItem.Invoke()).FirstOrDefault();
        var yObject = new List<Master>(yGroupItem.Invoke()).FirstOrDefault();
        if (xObject.DateTime > yObject.DateTime)
        {
            return -1;
        }
        else
        {
            return 1;
        }
    }
}